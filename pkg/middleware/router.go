package middleware

import (
	"github.com/coretrix/hitrix/pkg/controller"
	"github.com/gin-gonic/gin"
)

func Router(ginEngine *gin.Engine) {
	var dev *controller.DevPanelController
	{
		devGroup := ginEngine.Group("/dev/")
		//devGroup.Use(AuthorizeDevUser())
		{
			devGroup.GET("clear-cache/", dev.GetClearCacheAction)
			//devGroup.GET("clear-redis-streams/", dev.GetClearRedisStreamsAction)
			//devGroup.DELETE("delete-redis-stream/:name/", dev.DeleteRedisStreamAction)
			devGroup.GET("alters/", dev.GetAlters)
			devGroup.GET("redis-streams/", dev.GetRedisStreams)
			devGroup.GET("redis-statistics/", dev.GetRedisStatistics)

			devGroup.GET("login/", dev.PostLoginDevPanelAction)
			devGroup.GET("generate-token/", dev.PostGenerateTokenAction)
		}
	}
}
