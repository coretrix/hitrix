package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	graphqlParser "github.com/coretrix/hitrix/pkg/test/graphql-parser"
	"github.com/h2non/filetype"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
)

var dbAlters string
var parallelTestID string

type Environment struct {
	t                *testing.T
	Hitrix           *hitrix.Hitrix
	GinEngine        *gin.Engine
	Cxt              *gin.Context
	ResponseRecorder *httptest.ResponseRecorder
}

func (env *Environment) HandleQuery(query interface{}, variables map[string]interface{}, headers map[string]string) *graphqlParser.Errors {
	buff, err := graphqlParser.NewQueryParser().ParseQuery(query, variables)
	if err != nil {
		env.t.Fatal(err)
	}

	return env.handle(buff, query, headers)
}

func (env *Environment) HandleMutationMultiPart(query string, variables map[string]interface{},
	files map[string]string, variablesMap map[string]interface{}, headers map[string]string, model interface{}) (*graphqlParser.Errors, json.RawMessage) {
	return env.handleMultiPart(query, variables, files, variablesMap, headers, model)
}

func (env *Environment) HandleMutation(mutation interface{}, variables map[string]interface{}, headers map[string]string) *graphqlParser.Errors {
	buff, err := graphqlParser.NewQueryParser().ParseMutation(mutation, variables)
	if err != nil {
		env.t.Fatal(err)
	}

	return env.handle(buff, mutation, headers)
}

func (env *Environment) handle(buff bytes.Buffer, v interface{}, headers map[string]string) *graphqlParser.Errors {
	r, _ := http.NewRequestWithContext(env.Cxt, http.MethodPost, "/query", &buff)
	r.Header = http.Header{"Content-Type": []string{"application/json"}}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	env.Cxt.Request = r
	env.GinEngine.HandleContext(env.Cxt)

	var out struct {
		Data   *json.RawMessage
		Errors *graphqlParser.Errors
	}
	if err := json.NewDecoder(env.ResponseRecorder.Body).Decode(&out); err != nil {
		env.t.Fatal(err)
	}

	if out.Errors != nil {
		return out.Errors
	}

	if out.Data != nil {
		if err := json.Unmarshal(*out.Data, v); err != nil {
			env.t.Fatal(err)
		}
	}

	return nil
}
func (env *Environment) handleMultiPart(query string, variables interface{},
	files map[string]string, variablesMap map[string]interface{}, headers map[string]string, model interface{}) (*graphqlParser.Errors, json.RawMessage) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for name, path := range files {
		buf, _ := ioutil.ReadFile(path)
		kind, _ := filetype.Match(buf)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, name, path))
		h.Set("Content-Type", kind.MIME.Value)
		re, err := w.CreatePart(h)

		if err != nil {
			env.t.Fatal(err)
		}

		file, err := os.Open(path)

		defer func() {
			_ = file.Close()
		}()
		if err != nil {
			env.t.Fatal(err)
		}

		_, err = io.Copy(re, file)
		if err != nil {
			env.t.Fatal(err)
		}
	}

	jsonMap, err := json.Marshal(variablesMap)
	if err != nil {
		env.t.Fatal(err)
	}

	err = w.WriteField("map", string(jsonMap))
	if err != nil {
		env.t.Fatal(err)
	}

	queryJSON, _ := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	err = w.WriteField("operations", string(queryJSON))
	if err != nil {
		env.t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		env.t.Fatal(err)
	}

	r, _ := http.NewRequestWithContext(env.Cxt, http.MethodPost, "/query", &b)
	r.Header = http.Header{"Content-Type": []string{w.FormDataContentType()}}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	env.Cxt.Request = r
	env.GinEngine.HandleContext(env.Cxt)

	var out struct {
		Data   *json.RawMessage
		Errors *graphqlParser.Errors
	}
	if err := json.NewDecoder(env.ResponseRecorder.Body).Decode(&out); err != nil {
		env.t.Fatal(err)
	}

	if out.Errors != nil {
		return out.Errors, nil
	}

	return nil, *out.Data
}

func CreateContext(t *testing.T, projectName string, defaultServices []*service.DefinitionGlobal, mockServices ...*service.DefinitionGlobal) *Environment {
	var deferFunc func()

	err := os.Setenv("TZ", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("APP_MODE", app.ModeTest)
	if err != nil {
		t.Fatal(err)
	}

	testSpringInstance, deferFunc := hitrix.New(projectName, "").
		SetParallelTestID(getParallelID()).
		RegisterDIGlobalService(append(defaultServices, mockServices...)...).Build()
	defer deferFunc()

	var has bool
	ormService, has := service.DI().OrmEngine()
	if !has {
		panic("ORM is not loaded")
	}

	executeAlters(ormService)

	return &Environment{t: t, Hitrix: testSpringInstance}
}

func CreateAPIContext(t *testing.T, projectName string, resolvers graphql.ExecutableSchema, ginInitHandler hitrix.GinInitHandler, defaultGlobalServices []*service.DefinitionGlobal, defaultRequestServices []*service.DefinitionRequest, mockGlobalServices []*service.DefinitionGlobal, mockRequestServices []*service.DefinitionRequest) *Environment {
	var deferFunc func()
	gqlServerInitHandler := func(server *handler.Server) {
		server.AddTransport(transport.MultipartForm{})
	}
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("APP_MODE", app.ModeTest)
	if err != nil {
		t.Fatal(err)
	}

	testSpringInstance, deferFunc := hitrix.New(projectName, "").
		SetParallelTestID(getParallelID()).
		RegisterDIGlobalService(append(defaultGlobalServices, mockGlobalServices...)...).
		RegisterDIRequestService(append(defaultRequestServices, mockRequestServices...)...).Build()

	defer deferFunc()
	ginTestInstance := hitrix.InitGin(resolvers, ginInitHandler, gqlServerInitHandler)

	var has bool
	ormService, has := service.DI().OrmEngine()
	if !has {
		panic("ORM is not loaded")
	}

	executeAlters(ormService)

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	return &Environment{t: t, Hitrix: testSpringInstance, GinEngine: ginTestInstance, Cxt: c, ResponseRecorder: resp}
}

func executeAlters(ormService *beeorm.Engine) {
	if dbAlters == "" {
		err := dropTables(ormService.GetMysql())
		if err != nil {
			panic(err)
		}

		for _, alter := range ormService.GetAlters() {
			dbAlters += alter.SQL
		}

		_, def := ormService.GetMysql().Query(dbAlters)
		defer def()
	} else {
		err := truncateTables(ormService.GetMysql())
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		ormService.GetLocalCache().Clear()
		ormService.GetRedis().FlushAll()
	}

	altersSearch := ormService.GetRedisSearchIndexAlters()
	for _, alter := range altersSearch {
		alter.Execute()
	}
}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x%d", b, os.Getpid())[:10]
}

func getParallelID() string {
	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		return "1"
	} else if parallelTestID != "" {
		return parallelTestID
	} else {
		parallelTestID = getRandomString()
		return parallelTestID
	}
}

func dropTables(dbService *beeorm.DB) error {
	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('DROP TABLE ',table_schema,'.',table_name,';') AS query " +
			"FROM information_schema.tables WHERE table_schema IN ('" + dbService.GetPoolConfig().GetDatabase() + "')",
	)
	defer deferF()

	if rows != nil {
		var queries string

		for rows.Next() {
			rows.Scan(&query)
			queries += query
		}
		_, def := dbService.Query("SET FOREIGN_KEY_CHECKS=0;" + queries + "SET FOREIGN_KEY_CHECKS=1")

		defer def()
	}

	return nil
}

func truncateTables(dbService *beeorm.DB) error {
	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('delete from  ',table_schema,'.',table_name,';' , 'ALTER TABLE ', table_schema,'.',table_name , ' AUTO_INCREMENT = 1;') AS query " +
			"FROM information_schema.tables WHERE table_schema IN ('" + dbService.GetPoolConfig().GetDatabase() + "');",
	)
	defer deferF()

	if rows != nil {
		var queries string

		for rows.Next() {

			rows.Scan(&query)
			queries += query
		}

		_, def := dbService.Query("SET FOREIGN_KEY_CHECKS=0;" + queries + "SET FOREIGN_KEY_CHECKS=1")
		defer def()
	}

	return nil
}
