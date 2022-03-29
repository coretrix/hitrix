package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

func Cors(ginEngine *gin.Engine) {
	configService := service.DI().Config()
	origins, ok := configService.Strings("cors")
	if !ok {
		panic("cors is missing")
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		AllowOrigins:     origins,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"X-Invalid-Authorization", "X-Clockwork-Id", "X-Clockwork-Version"},
	}

	ginEngine.Use(cors.New(corsConfig))
}
