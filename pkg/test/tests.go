package test

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	"github.com/summer-solutions/orm"
)

var dbService *orm.DB
var ormService *orm.Engine
var ginTestInstance *gin.Engine
var testSpringInstance *hitrix.Hitrix

type Ctx struct {
	t *testing.T
	g *gin.Engine
	c *gin.Context
	w *httptest.ResponseRecorder
}

func (ctx *Ctx) HandleQuery(query interface{}, variables map[string]interface{}) *graphqlParser.Errors {
	buff, err := graphqlParser.NewQueryParser().ParseQuery(query, variables)
	if err != nil {
		ctx.t.Fatal(err)
	}

	return ctx.handle(buff, query)
}

func (ctx *Ctx) HandleMutation(mutation interface{}, variables map[string]interface{}) *graphqlParser.Errors {
	buff, err := graphqlParser.NewQueryParser().ParseMutation(mutation, variables)
	if err != nil {
		ctx.t.Fatal(err)
	}

	return ctx.handle(buff, mutation)
}

func (ctx *Ctx) handle(buff bytes.Buffer, v interface{}) *graphqlParser.Errors {
	r, _ := http.NewRequestWithContext(ctx.c, http.MethodPost, "/query", &buff)
	r.Header = http.Header{"Content-Type": []string{"application/json"}}
	ctx.c.Request = r
	ctx.g.HandleContext(ctx.c)

	var out struct {
		Data   *json.RawMessage
		Errors *graphqlParser.Errors
	}
	if err := json.NewDecoder(ctx.w.Body).Decode(&out); err != nil {
		ctx.t.Fatal(err)
	}

	if out.Errors != nil {
		return out.Errors
	}

	if out.Data != nil {
		if err := jsonutil.UnmarshalGraphQL(*out.Data, v); err != nil {
			ctx.t.Fatal(err)
		}
	}

	return nil
}

func CreateContext(t *testing.T, projectName string, resolvers graphql.ExecutableSchema, defaultServices []*service.Definition, mockServices ...*service.Definition) *Ctx {
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
		ginTestInstance = hitrix.InitGin(resolvers, nil)

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
		ginTestInstance = hitrix.InitGin(resolvers, nil)

		// TODO: fix multiple connections to mysql
		ormService, _ = service.DI().OrmEngine()
		dbService = ormService.GetMysql()
	}

	err := truncateTables()
	if err != nil {
		t.Fatal(err)
	}

	ormService.GetLocalCache().Clear()
	ormService.GetRedis().FlushDB()

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	return &Ctx{t: t, g: ginTestInstance, c: c, w: resp}
}

func dropTables() error {
	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('DROP TABLE ',table_schema,'.',table_name,';') AS query " +
			"FROM information_schema.tables WHERE table_schema IN ('" + dbService.GetDatabaseName() + "_test')",
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
			"FROM information_schema.tables WHERE table_schema IN ('" + dbService.GetDatabaseName() + "_test');",
	)
	defer deferF()
	fmt.Println(dbService.GetDatabaseName())
	fmt.Println(rows == nil)
	if rows != nil {
		fmt.Println("query")
		var queries string

		for rows.Next() {
			rows.Scan(&query)
			fmt.Println(query)
			queries += query
		}
		fmt.Println(queries)

		_, def := dbService.Query("SET FOREIGN_KEY_CHECKS=0;" + queries + "SET FOREIGN_KEY_CHECKS=1")
		defer def()
	}

	return nil
}
