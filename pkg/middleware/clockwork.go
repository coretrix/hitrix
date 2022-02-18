package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func (h *clockWorkHandler) Handle(log map[string]interface{}) {
	if log["source"] == "mysql" {
		microseconds := log["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)
		query := log["query"].(string)
		h.DatabaseDataSource.LogQuery("mysql", query, milliseconds, nil)
	} else if log["source"] == "redis" {
		_, hasMisses := log["miss"]
		operation := log["operation"].(string)
		if strings.Contains(operation, "profiler_store.") {
			return
		}

		microseconds := log["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)

		if hasMisses {
			h.RedisDataSource.LogCacheMiss(log["pool"].(string), operation, log["query"].(string), "", 1, milliseconds, 0)
		} else {
			h.RedisDataSource.LogCache(log["pool"].(string), dataSource.CacheHit, operation, log["query"].(string), "", milliseconds, 0)
		}
	} else if log["source"] == "elastic" {
		microseconds := log["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)

		microsecondsQueryTime := log["query_time"].(int64)
		millisecondsQueryTime := float32(microsecondsQueryTime) / float32(1000)

		query := log["post"]
		if sort, ok := log["sort"]; ok {
			query = fmt.Sprint(query, fmt.Sprintf(" SORT: %s", sort))
		}
		query = fmt.Sprint(query, fmt.Sprintf(" QUERY TIME:%.3f ms", millisecondsQueryTime))
		h.DatabaseDataSource.LogQuery("elastic", query.(string), milliseconds, nil)
	} else if log["source"] == "clickhouse" {
		microseconds := log["microseconds"].(int64)
		milliseconds := float32(microseconds) / float32(1000)
		query := fmt.Sprint(log["query"])
		h.DatabaseDataSource.LogQuery("clickhouse", query, milliseconds, nil)
	} else if log["source"] == "local_cache" {
		query := strings.Split(log["query"].(string), " ")
		key := strings.Split(query[2], ":")
		tableSchema := h.ormService.GetRegistry().GetTableSchemaForCachePrefix(key[0])

		h.LocalCacheDataSource.LogTable(
			map[string]interface{}{
				"Operation": log["operation"],
				"Query":     tableSchema.GetTableName() + ":" + key[1] + " " + query[3],
			}, "Queries", nil)
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

		clockWorkHandler := clockWorkHandler{ormService: ormService, DatabaseDataSource: databaseDataSource, RedisDataSource: redisDataSource, LocalCacheDataSource: localCacheDataSource}
		ormService.RegisterQueryLogger(&clockWorkHandler, true, true, true)

		profilerKey := c.Request.Header.Get("CoreTrix")

		if profilerKey != password {
			return
		}

		profilerService.GetRequestDataSource().SetStartTime(time.Now())
		profilerService.GetRequestDataSource().StartMemoryUsage()

		c.Writer.Header().Set("X-Clockwork-Id", profilerService.GetUniqueId())
		c.Writer.Header().Set("X-Clockwork-Version", "5.1.0")

		c.Next()

		profilerService.GetRequestDataSource().SetMiddleware(c.HandlerNames())
		profilerService.GetRequestDataSource().SetResponseTime(time.Now())
		profilerService.GetRequestDataSource().EndMemoryUsage()

		profilerService.GetRequestDataSource().SetResponseStatus(int16(c.Writer.Status()))

		profilerService.SaveData()
	})
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
