package test

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
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

func CreateContext(
	t *testing.T,
	projectName string,
	defaultServices []*service.DefinitionGlobal,
	mockGlobalServices []*service.DefinitionGlobal,
	redisPools *app.RedisPools,
) *Environment {
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
		RegisterDIGlobalService(append(defaultServices, mockGlobalServices...)...).RegisterRedisPools(
		redisPools,
	).Build()
	defer deferFunc()

	ormService := service.DI().OrmEngine()

	executeAlters(ormService)

	return &Environment{t: t, Hitrix: testSpringInstance}
}

func CreateAPIContext(
	t *testing.T,
	projectName string,
	ginInitHandler hitrix.GinInitHandler,
	defaultGlobalServices []*service.DefinitionGlobal,
	defaultRequestServices []*service.DefinitionRequest,
	mockGlobalServices []*service.DefinitionGlobal,
	mockRequestServices []*service.DefinitionRequest,
	redisPools *app.RedisPools,
) *Environment {
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
		RegisterDIGlobalService(append(defaultGlobalServices, mockGlobalServices...)...).
		RegisterDIRequestService(append(defaultRequestServices, mockRequestServices...)...).
		RegisterRedisPools(
			redisPools,
		).Build()

	defer deferFunc()

	ginTestInstance := hitrix.InitGin(ginInitHandler)

	ormService := service.DI().OrmEngine()

	executeAlters(ormService)

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)

	return &Environment{t: t, Hitrix: testSpringInstance, GinEngine: ginTestInstance, Cxt: c, ResponseRecorder: resp}
}

func executeAlters(ormService *beeorm.Engine) {
	if dbAlters == "" {
		dropTables(ormService.GetMysql())

		for _, alter := range ormService.GetAlters() {
			dbAlters += alter.SQL
		}

		_, def := ormService.GetMysql().Query(dbAlters)
		defer def()
	} else {
		truncateTables(ormService.GetMysql())
	}

	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		ormService.GetLocalCache().Clear()

		pools := service.DI().App().RedisPools
		ormService.GetRedis(pools.Stream).FlushDB()
		ormService.GetRedis(pools.Persistent).FlushDB()
		ormService.GetRedis(pools.Cache).FlushDB()

		for _, pool := range pools.Search {
			ormService.GetRedis(pool).FlushDB()
		}
	}

	altersSearch := ormService.GetRedisSearchIndexAlters()
	for _, alter := range altersSearch {
		alter.Execute()
	}
}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 10)

	//nolint //G404: Use of weak random number generator (math/rand instead of crypto/rand)
	rand.Read(b)

	return fmt.Sprintf("%x%d", b, os.Getpid())[:5]
}

func getParallelID() string {
	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		return "1"
	} else if parallelTestID != "" {
		return parallelTestID
	}

	parallelTestID = getRandomString()

	return parallelTestID
}

func dropTables(dbService *beeorm.DB) {
	var query string
	rows, deferF := dbService.Query(
		"SELECT CONCAT('DROP TABLE IF EXISTS ',table_schema,'.',table_name,';') AS query " +
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
}

func truncateTables(dbService *beeorm.DB) {
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
}
