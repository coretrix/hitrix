package test

import (
	"bytes"
	"encoding/json"

	graphqlParser "github.com/coretrix/hitrix/pkg/test/graphql-parser"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coretrix/hitrix/pkg/test/graphql-parser/jsonutil"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/orm"
)

var dbService *orm.DB
var ormService *orm.Engine
var ginTestInstance *gin.Engine
var testSpringInstance *hitrix.Hitrix

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
		if err := jsonutil.UnmarshalGraphQL(*out.Data, v); err != nil {
			env.t.Fatal(err)
		}
	}

	return nil
}

func CreateContext(t *testing.T, projectName string, defaultServices []*service.Definition, mockServices ...*service.Definition) *Environment {
	var deferFunc func()

	if testSpringInstance == nil {
		err := os.Setenv("TZ", "UTC")
		if err != nil {
			t.Fatal(err)
		}
		err = os.Setenv("APP_MODE", app.ModeTest)
		if err != nil {
			t.Fatal(err)
		}

		testSpringInstance, deferFunc = hitrix.New(projectName, "").RegisterDIService(append(defaultServices, mockServices...)...).Build()
		defer deferFunc()

		var has bool
		ormService, has = service.DI().OrmEngine()
		if !has {
			panic("ORM is not loaded")
		}

		dbService = ormService.GetMysql()

		err = dropTables()
		if err != nil {
			t.Fatal(err)
		}

		alters := ormService.GetAlters()

		var queries string

		for _, alter := range alters {
			queries += alter.SQL
		}

		if queries != "" {
			_, def := dbService.Query(queries)
			defer def()
		}
	}

	if len(mockServices) != 0 {
		testSpringInstance, deferFunc = hitrix.New(projectName, "").RegisterDIService(append(defaultServices, mockServices...)...).Build()
		defer deferFunc()

		// TODO: fix multiple connections to mysql
		ormService, _ = service.DI().OrmEngine()
		dbService = ormService.GetMysql()
	}

	err := truncateTables()
	if err != nil {
		t.Fatal(err)
	}

	ormService.GetLocalCache().Clear()
	ormService.GetRedis().FlushAll()

	altersSearch := ormService.GetRedisSearchIndexAlters()
	for _, alter := range altersSearch {
		alter.Execute()
	}

	return &Environment{t: t, Hitrix: testSpringInstance, GinEngine: ginTestInstance}
}

func CreateAPIContext(t *testing.T, projectName string, resolvers graphql.ExecutableSchema, ginInitHandler hitrix.GinInitHandler, defaultServices []*service.Definition, mockServices ...*service.Definition) *Environment {
	var deferFunc func()

	if testSpringInstance == nil {
		err := os.Setenv("TZ", "UTC")
		if err != nil {
			t.Fatal(err)
		}
		err = os.Setenv("APP_MODE", app.ModeTest)
		if err != nil {
			t.Fatal(err)
		}

		testSpringInstance, deferFunc = hitrix.New(projectName, "").RegisterDIService(append(defaultServices, mockServices...)...).Build()
		defer deferFunc()
		ginTestInstance = hitrix.InitGin(resolvers, ginInitHandler, nil)

		var has bool
		ormService, has = service.DI().OrmEngine()
		if !has {
			panic("ORM is not loaded")
		}

		dbService = ormService.GetMysql()

		err = dropTables()
		if err != nil {
			t.Fatal(err)
		}

		alters := ormService.GetAlters()

		var queries string

		for _, alter := range alters {
			queries += alter.SQL
		}

		if queries != "" {
			_, def := dbService.Query(queries)
			defer def()
		}
	}

	if len(mockServices) != 0 {
		testSpringInstance, deferFunc = hitrix.New(projectName, "").RegisterDIService(append(defaultServices, mockServices...)...).Build()
		defer deferFunc()
		ginTestInstance = hitrix.InitGin(resolvers, ginInitHandler, nil)

		// TODO: fix multiple connections to mysql
		ormService, _ = service.DI().OrmEngine()
		dbService = ormService.GetMysql()
	}

	err := truncateTables()
	if err != nil {
		t.Fatal(err)
	}

	ormService.GetLocalCache().Clear()
	ormService.GetRedis().FlushAll()

	altersSearch := ormService.GetRedisSearchIndexAlters()
	for _, alter := range altersSearch {
		alter.Execute()
	}

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	return &Environment{t: t, Hitrix: testSpringInstance, GinEngine: ginTestInstance, Cxt: c, ResponseRecorder: resp}
}

func dropTables() error {
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

func truncateTables() error {
	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('truncate table ',table_schema,'.',table_name,';') AS query " +
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
