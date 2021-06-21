package middleware

import (
	"time"

	"github.com/coretrix/hitrix/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(ginEngine *gin.Engine) {
	configService := service.DI().Config()
	origins, ok := configService.Strings("cors")
	if !ok {
		panic("cors is missing")
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOrigins:     origins,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"X-Invalid-Authorization"},
	}

	ginEngine.Use(cors.New(corsConfig))
}
