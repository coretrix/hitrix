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

	var dev *controller.DevPanelController
	{
		devGroup := ginEngine.Group("/dev/")
		devGroup.Use(AuthorizeDevUser())
		{
			devGroup.GET("clear-cache/", dev.GetClearCacheAction)
			//devGroup.GET("clear-redis-streams/", dev.GetClearRedisStreamsAction)
			//devGroup.DELETE("delete-redis-stream/:name/", dev.DeleteRedisStreamAction)
			devGroup.GET("alters/", dev.GetAlters)
			devGroup.GET("redis-streams/", dev.GetRedisStreams)
			devGroup.GET("redis-statistics/", dev.GetRedisStatistics)

			devGroup.POST("login/", dev.PostLoginDevPanelAction)
			devGroup.POST("generate-token/", AuthorizeWithDevRefreshToken(), dev.PostGenerateTokenAction)
		}
	}
}
