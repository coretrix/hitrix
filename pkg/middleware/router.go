package middleware

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/pkg/controller"
	"github.com/gin-gonic/gin"
)

func Router(ginEngine *gin.Engine) {
	_, has := hitrix.DIC().JWT()
	if !has {
		panic("Please load JWT service")
	}

	var devPanel *controller.DevPanelController
	{
		devGroup := ginEngine.Group("/dev/")
		devGroup.Use(AuthorizeDevUser())
		{
			ginEngine.GET("/dev/action-list/", devPanel.GetActionListAction) //todo for remove

			devGroup.GET("clear-cache/", devPanel.GetClearCacheAction)
			devGroup.GET("clear-redis-streams/", devPanel.GetClearRedisStreamsAction)
			devGroup.DELETE("delete-redis-stream/:name/", devPanel.DeleteRedisStreamAction)
			devGroup.GET("alters/", devPanel.GetAlters)
			devGroup.GET("redis-streams/", devPanel.GetRedisStreams)
			devGroup.GET("redis-statistics/", devPanel.GetRedisStatistics)

			ginEngine.GET("dev/create-admin/", devPanel.CreateAdminUserAction)
			ginEngine.POST("dev/login/", devPanel.PostLoginDevPanelAction)
			ginEngine.POST("dev/generate-token/", AuthorizeWithDevRefreshToken(), devPanel.PostGenerateTokenAction)
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
