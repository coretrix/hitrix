package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/controller"
	"github.com/coretrix/hitrix/service"
)

func Router(ginEngine *gin.Engine) {
	_ = service.DI().JWT()

	var devPanel *controller.DevPanelController
	{
		devGroup := ginEngine.Group("/dev/")
		if !service.DI().App().IsInLocalMode() {
			devGroup.Use(AuthorizeDevUser())
		}
		{
			//ginEngine.GET("/dev/action-list/", devPanel.GetActionListAction)

			devGroup.GET("settings/", devPanel.GetSettingsAction)
			devGroup.GET("clear-cache/", devPanel.GetClearCacheAction)
			devGroup.GET("clear-redis-streams/", devPanel.GetClearRedisStreamsAction)
			devGroup.DELETE("delete-redis-stream/:name/", devPanel.DeleteRedisStreamAction)
			devGroup.GET("alters/", devPanel.GetAlters)
			devGroup.GET("redis-search/alters/", devPanel.GetRedisSearchAlters)
			devGroup.GET("redis-search/statistics/", devPanel.GetRedisSearchStatistics)
			devGroup.GET("redis-search/indexes/", devPanel.GetRedisSearchIndexes)
			devGroup.GET("redis-search/index/info/:index/", devPanel.PostRedisSearchIndexInfo)
			devGroup.GET("redis-search/force-reindex/:index/", devPanel.PostRedisSearchForceReindex)

			devGroup.GET("redis-streams/", devPanel.GetRedisStreams)
			devGroup.GET("redis-statistics/", devPanel.GetRedisStatistics)

			ginEngine.GET("dev/create-dev-panel-user/", devPanel.CreateDevPanelUserAction)
			ginEngine.POST("dev/login/", devPanel.PostLoginDevPanelAction)
			ginEngine.POST("dev/generate-token/", AuthorizeWithDevRefreshToken(), devPanel.PostGenerateTokenAction)

			devGroup.GET("feature-flag/list/", devPanel.GetFeatureFlags)
			devGroup.POST("feature-flag/enable/:name/", devPanel.PostEnableFeatureFlag)
			devGroup.POST("feature-flag/disable/:name/", devPanel.PostDisableFeatureFlag)
		}
	}

	var errorLog *controller.ErrorLogController
	{
		errorLogGroup := ginEngine.Group("/error-log/")
		errorLogGroup.Use(AuthorizeDevUser())

		errorLogGroup.GET("errors/", errorLog.GetErrors)
		errorLogGroup.GET("remove/:id/", errorLog.DeleteError)
		errorLogGroup.GET("remove-all/", errorLog.DeleteAllErrors)
		errorLogGroup.GET("panic/", errorLog.Panic)
	}
}
