package controller

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"
	"github.com/latolukasz/beeorm/tools"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/indexes"
	"github.com/coretrix/hitrix/pkg/dto/list"
	"github.com/coretrix/hitrix/pkg/entity"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	"github.com/coretrix/hitrix/pkg/errors"
	accountModel "github.com/coretrix/hitrix/pkg/model/account"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/coretrix/hitrix/pkg/view/metrics"
	"github.com/coretrix/hitrix/pkg/view/requestlogger"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type MenuItem struct {
	Label string
	URL   string
	Icon  string
}

type DevPanelController struct {
}

//func (controller *DevPanelController) GetActionListAction(c *gin.Context) {
//	actions := []*MenuItem{
//		{
//			Label: "Clear Cache",
//			URL:   "/dev/clear-cache/",
//			Icon:  "mdiCached",
//		},
//	}
//
//	c.JSON(200, actions)
//}

func (controller *DevPanelController) GetSettingsAction(c *gin.Context) {
	appService := service.DI().App()

	response.SuccessResponse(c, gin.H{
		"AppMode": appService.Mode,
	})
}

func (controller *DevPanelController) CreateDevPanelUserAction(c *gin.Context) {
	passwordService := service.DI().Password()

	ormService := service.DI().OrmEngine()

	form := &accountModel.LoginDevForm{}
	if err := binding.ShouldBindQuery(c, form); err != nil {
		fieldError, ok := (err).(errors.FieldErrors)
		if ok {
			response.ErrorResponseFields(c, fieldError, nil)

			return
		}

		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	adminEntity := service.DI().App().DevPanel.UserEntity

	passwordHash, err := passwordService.HashPassword(form.Password)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
	}

	adminTableSchema := ormService.GetRegistry().GetTableSchemaForEntity(adminEntity)
	response.SuccessResponse(
		c,
		fmt.Sprintf(`INSERT INTO %s (Email, Password) VALUES('%s', '%s')`, adminTableSchema.GetTableName(), form.Username, passwordHash))
}

func (controller *DevPanelController) PostLoginDevPanelAction(c *gin.Context) {
	loginForm := accountModel.LoginDevForm{}
	token, refreshToken, err := loginForm.Login(c)

	errType, ok := err.(errors.FieldErrors)

	if ok && errType != nil {
		response.ErrorResponseFields(c, errType, nil)

		return
	}

	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) PostGenerateTokenAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	devPanelUserEntity := c.MustGet(account.LoggedDevPanelUserEntity).(app.IDevPanelUserEntity)

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) GetClearCacheAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	redisService := ormService.GetRedis()

	redisService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetClearRedisStreamsAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Stream == "" {
		panic("stream pool is not defined")
	}

	redisStreamsService := ormService.GetRedis(appService.RedisPools.Stream)
	redisStreamsService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) DeleteRedisStreamAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Stream == "" {
		panic("stream pool is not defined")
	}

	redisStreamService := ormService.GetRedis(appService.RedisPools.Stream)

	name := c.Param("name")
	if name == "" {
		panic("provide stream name")
	}

	redisStreamService.XTrim(name, 0)

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetAlters(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	alters := ormService.GetAlters()
	result := make([]string, len(alters))

	force := c.Query("force")
	if force != "" {
		redisService := ormService.GetRedis()
		redisService.FlushDB()
	}

	for i, alter := range alters {
		if force != "" {
			alter.Exec()
		} else {
			result[i] = alter.SQL
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisStreams(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	stats := tools.GetRedisStreamsStatistics(ormService)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Stream < stats[j].Stream
	})
	response.SuccessResponse(c, stats)
}

func (controller *DevPanelController) GetRedisStatistics(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	stats := tools.GetRedisStatistics(ormService)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].RedisPool < stats[j].RedisPool
	})
	response.SuccessResponse(c, stats)
}

func (controller *DevPanelController) GetRedisSearchStatistics(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	response.SuccessResponse(c, tools.GetRedisSearchStatistics(ormService))
}

func (controller *DevPanelController) GetRedisSearchAlters(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	altersSearch := ormService.GetRedisSearchIndexAlters()
	result := make([]map[string]string, len(altersSearch))

	force := c.Query("force")
	for i, alter := range altersSearch {
		if force != "" {
			alter.Execute()
		} else {
			result[i] = map[string]string{
				"Query":   alter.Query,
				"Changes": strings.Join(alter.Changes, " | "),
			}
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisSearchIndexes(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	indices := ormService.GetRedisSearch(appService.RedisPools.Search).ListIndices()
	sort.Strings(indices)

	indexList := make([]indexes.Index, len(indices))

	for i, indexName := range indices {
		info := ormService.GetRedisSearch(appService.RedisPools.Search).Info(indexName)
		indexList[i] = indexes.Index{
			Name:      indexName,
			TotalDocs: info.NumDocs,
			TotalSize: uint64(info.DocTableSizeMB + info.KeyTableSizeMB + info.SortableValuesSizeMB + info.InvertedSzMB + info.OffsetVectorsSzMB),
		}
	}

	result := indexes.ResponseDTOList{
		Indexes: indexList,
	}
	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) PostRedisSearchForceReindex(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	indexName := c.Param("index")
	if indexName == "" {
		response.ErrorResponseGlobal(c, "index is required", nil)

		return
	}

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	ormService.GetRedisSearch(appService.RedisPools.Search).ForceReindex(indexName)
	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostRedisSearchForceReindexAll(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	indexes := ormService.GetRedisSearch(appService.RedisPools.Search).ListIndices()

	concurrently := c.Query("concurrently")
	if concurrently != "" {
		redisSearch := ormService.GetRedisSearch(appService.RedisPools.Search)

		wg := sync.WaitGroup{}
		wg.Add(len(indexes))

		for _, index := range indexes {
			go func(index string) {
				defer func() {
					if r := recover(); r != nil {
						service.DI().ErrorLogger().LogError(r)
					}
				}()

				redisSearch.ForceReindex(index)
				wg.Done()
			}(index)
		}

		wg.Wait()
	} else {
		for _, index := range indexes {
			ormService.GetRedisSearch(appService.RedisPools.Search).ForceReindex(index)
		}
	}

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostRedisSearchIndexInfo(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	indexName := c.Param("index")
	if indexName == "" {
		response.ErrorResponseGlobal(c, "index is required", nil)

		return
	}

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	response.SuccessResponse(c, ormService.GetRedisSearch(appService.RedisPools.Search).Info(indexName))
}

func (controller *DevPanelController) GetFeatureFlags(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, beeorm.NewPager(1, 1000))

	type feature struct {
		Name       string
		Registered bool
		Enabled    bool
	}

	result := make([]*feature, len(featureFlagEntities))

	for i, featureFlagEntity := range featureFlagEntities {
		result[i] = &feature{
			Name:       featureFlagEntity.Name,
			Registered: featureFlagEntity.Registered,
			Enabled:    featureFlagEntity.Enabled,
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) PostEnableFeatureFlag(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.ErrorResponseGlobal(c, "name is required", nil)

		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)

	if !found {
		response.ErrorResponseGlobal(c, "feature is missing", nil)

		return
	}

	featureFlagEntity.Enabled = true
	ormService.Flush(featureFlagEntity)

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostDisableFeatureFlag(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.ErrorResponseGlobal(c, "name is required", nil)

		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)

	if !found {
		response.ErrorResponseGlobal(c, "feature is missing", nil)

		return
	}

	featureFlagEntity.Enabled = false
	ormService.Flush(featureFlagEntity)

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) GetEnvValues(c *gin.Context) {
	response.SuccessResponse(c, map[string]interface{}{"list": os.Environ()})
}

func (controller *DevPanelController) PostRequestsLogger(c *gin.Context) {
	request := list.RequestDTOList{}

	err := binding.ShouldBindJSON(c, &request)
	if errorhandling.HandleError(c, err) {
		return
	}

	res, err := requestlogger.RequestsLogger(c.Request.Context(), request)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, res)
}

func (controller *DevPanelController) GetMetrics(c *gin.Context) {
	res := metrics.Get(c.Request.Context())

	response.SuccessResponse(c, res)
}
