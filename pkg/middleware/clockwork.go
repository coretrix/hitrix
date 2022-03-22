package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/anton-shumanski/clockwork"
	dataSource "github.com/anton-shumanski/clockwork/data-source"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
)

type clockWorkHandler struct {
	ormService           *beeorm.Engine
	DatabaseDataSource   dataSource.QueryLoggerDataSourceInterface
	RedisDataSource      dataSource.CacheLoggerDataSourceInterface
	LocalCacheDataSource dataSource.UserDataSourceInterface
}

func (h *clockWorkHandler) Handle(logData map[string]interface{}) {
	if logData["source"] == "mysql" {
		microseconds := logData["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)
		query := logData["query"].(string)
		h.DatabaseDataSource.LogQuery("mysql", query, milliseconds, nil)
	} else if logData["source"] == "redis" {
		_, hasMisses := logData["miss"]
		operation := logData["operation"].(string)
		if strings.Contains(operation, "profiler_store.") {
			return
		}

		microseconds := logData["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)

		if hasMisses {
			h.RedisDataSource.LogCacheMiss(logData["pool"].(string), operation, logData["query"].(string), "", 1, milliseconds, 0)
		} else {
			h.RedisDataSource.LogCache(logData["pool"].(string), dataSource.CacheHit, operation, logData["query"].(string), "", milliseconds, 0)
		}
	} else if logData["source"] == "elastic" {
		microseconds := logData["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)

		microsecondsQueryTime := logData["query_time"].(int64)
		millisecondsQueryTime := float32(microsecondsQueryTime) / float32(1000)

		query := logData["post"]
		if sort, ok := logData["sort"]; ok {
			query = fmt.Sprint(query, fmt.Sprintf(" SORT: %s", sort))
		}
		query = fmt.Sprint(query, fmt.Sprintf(" QUERY TIME:%.3f ms", millisecondsQueryTime))
		h.DatabaseDataSource.LogQuery("elastic", query.(string), milliseconds, nil)
	} else if logData["source"] == "clickhouse" {
		microseconds := logData["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)
		query := fmt.Sprint(logData["query"])
		h.DatabaseDataSource.LogQuery("clickhouse", query, milliseconds, nil)
	} else if logData["source"] == "local_cache" {
		originalQuery := logData["query"].(string)

		queryParts := strings.Fields(originalQuery)

		operator := queryParts[0]

		queries := strings.TrimSpace(strings.TrimPrefix(originalQuery, operator))

		var q string

		switch operator {
		case "MGET":
			originalKeys := strings.Split(queries, " ")

			for _, originalKey := range originalKeys {
				originalKeyArray := strings.Split(originalKey, ":")
				tableSchema := h.ormService.GetRegistry().GetTableSchemaForCachePrefix(originalKeyArray[0])

				q += tableSchema.GetTableName() + ":" + originalKeyArray[1] + " "
			}
		case "MSET":
			keyValues := strings.Split(queries, " ")
			for _, keyValue := range keyValues {
				tableAndIDPair := strings.Split(keyValue, ":")
				if len(tableAndIDPair) == 2 {
					tableSchema := h.ormService.GetRegistry().GetTableSchemaForCachePrefix(tableAndIDPair[0])
					q += tableSchema.GetTableName() + ":" + tableAndIDPair[1] + " "
				}
			}
		case "GET":
			originalKeyArray := strings.Split(queries, ":")
			tableSchema := h.ormService.GetRegistry().GetTableSchemaForCachePrefix(originalKeyArray[0])

			q += tableSchema.GetTableName() + ":" + originalKeyArray[1]

		case "SET":
			keyValue := strings.Split(queries, " ")
			originalKeyArray := strings.Split(keyValue[0], ":")
			tableSchema := h.ormService.GetRegistry().GetTableSchemaForCachePrefix(originalKeyArray[0])

			q += tableSchema.GetTableName() + ":" + originalKeyArray[1] + " "
		default:
			h.LocalCacheDataSource.LogTable(
				map[string]interface{}{
					"Operation": logData["operation"],
					"Query":     originalQuery,
				}, "Queries",
				nil,
			)
			return
		}

		h.LocalCacheDataSource.LogTable(
			map[string]interface{}{
				"Operation": logData["operation"],
				"Query":     q,
			}, "Queries",
			nil,
		)
	}
}

func Clockwork(ginEngine *gin.Engine) {
	ginEngine.Use(func(c *gin.Context) {
		configService := service.DI().Config()

		password, ok := configService.String("clockwork.password")
		if !ok {
			return
		}
		ormService := service.DI().OrmEngineForContext(c.Request.Context())

		redisDataProvider := &ormDataProvider{RedisStorageProvider: ormService.GetRedis()}
		profilerService := service.DI().ClockWorkForContext(c.Request.Context())
		profilerService.DataProvider = redisDataProvider

		var loggerDataSource dataSource.LoggerDataSourceInterface = &dataSource.LoggerDataSource{}
		profilerService.SetLoggerDataSource(loggerDataSource)

		var databaseDataSource dataSource.QueryLoggerDataSourceInterface = &dataSource.DatabaseDataSource{}
		profilerService.SetDatabaseDataSource(databaseDataSource)

		var redisDataSource dataSource.CacheLoggerDataSourceInterface = &dataSource.CacheDataSource{}
		profilerService.AddDataSource(redisDataSource)

		var requestDataSource dataSource.RequestLoggerDataSourceInterface = &dataSource.RequestResponseDataSource{}
		profilerService.SetRequestDataSource(requestDataSource)

		var localCacheDataSource dataSource.UserDataSourceInterface = new(dataSource.UserDataDataSource)
		profilerService.AddDataSource(localCacheDataSource)

		localCacheDataSource.SetShowAs("table")
		localCacheDataSource.SetTitle("Local Cache")

		if profilerKey := c.Request.Header.Get("CoreTrix"); profilerKey != password {
			return
		}

		if resolverName := setController(c, profilerService); resolverName == "IntrospectionQuery" {
			return
		}

		clockWorkHandler := clockWorkHandler{
			ormService:           ormService,
			DatabaseDataSource:   databaseDataSource,
			RedisDataSource:      redisDataSource,
			LocalCacheDataSource: localCacheDataSource,
		}

		ormService.RegisterQueryLogger(&clockWorkHandler, true, true, true)

		profilerService.GetRequestDataSource().SetStartTime(time.Now())
		profilerService.GetRequestDataSource().StartMemoryUsage()

		c.Writer.Header().Set("X-Clockwork-Id", profilerService.GetUniqueId())
		c.Writer.Header().Set("X-Clockwork-Version", "5.1.0")

		c.Next()

		middlewares := c.HandlerNames()
		profilerService.GetRequestDataSource().SetMiddleware(middlewares)
		profilerService.GetRequestDataSource().SetResponseTime(time.Now())
		profilerService.GetRequestDataSource().EndMemoryUsage()

		profilerService.GetRequestDataSource().SetResponseStatus(int16(c.Writer.Status()))

		profilerService.SaveData()
	})
}

func isMultipartRequest(c *gin.Context) bool {
	contentTypes := c.Request.Header["Content-Type"]
	return len(contentTypes) > 0 && strings.Contains(contentTypes[0], gin.MIMEMultipartPOSTForm)
}

func getRawData(c *gin.Context) []byte {
	var b []byte
	if isMultipartRequest(c) {
		err := c.Request.ParseMultipartForm(4096)
		if err != nil {
			panic(err)
		}
		v := c.Request.FormValue("operations")
		b = []byte(v)
	} else {
		var err error
		b, err = c.GetRawData()
		if err != nil {
			panic(err)
		}
	}
	return b
}

func setController(c *gin.Context, profilerService *clockwork.Clockwork) string {
	b := getRawData(c)

	if len(b) == 0 {
		return ""
	}

	bodyMap := map[string]interface{}{}
	err := json.Unmarshal(b, &bodyMap)
	if err != nil {
		panic(err)
	}
	var queryType string
	var resolverName string

	operationName, okOperationName := bodyMap["operationName"]
	queryValue, okQuery := bodyMap["query"]
	mutationValue, okMutation := bodyMap["mutation"]

	if okOperationName && (okQuery || okMutation) {
		if okQuery {
			queryType = "query"

			if operationName != nil {
				resolverName = operationName.(string)
			} else {
				resolverName = strings.TrimLeft(queryValue.(string), "{")
				resolverName = resolverName[:strings.IndexByte(resolverName, '{')]
			}
		} else if okMutation {
			queryType = "mutation"

			if operationName != nil {
				resolverName = operationName.(string)
			} else {
				resolverName = strings.TrimLeft(mutationValue.(string), "{")
				resolverName = resolverName[:strings.IndexByte(resolverName, '{')]
			}
		}
	}

	profilerService.GetRequestDataSource().SetController(queryType, resolverName)

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(b))

	return resolverName
}

type ormDataProvider struct {
	RedisStorageProvider *beeorm.RedisCache
}

func (provider *ormDataProvider) Get(key string, id string) dataSource.DataBuffer {
	value, has := provider.RedisStorageProvider.Get(key + "." + id)
	var raw dataSource.DataBuffer

	if !has {
		return raw
	}

	err := json.Unmarshal([]byte(value), &raw)
	if err != nil {
		service.DI().ErrorLogger().LogError(err)
	}

	return raw
}

func (provider *ormDataProvider) Set(key string, id string, data *dataSource.DataBuffer) {
	jsonString, _ := json.Marshal(data)
	provider.RedisStorageProvider.Set(key+"."+id, jsonString, int(time.Minute*5))
}
