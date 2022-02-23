package middleware

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/example/rest/controller"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"
)

func Router(ginEngine *gin.Engine) {
	var websocketController *controller.WebsocketController
	ginEngine.GET("/ws/", websocketController.InitConnection)

	ginEngine.GET("/dev/upload/", func(context *gin.Context) {
		object := service.DI().OSService().GetBucketConfig(oss.BucketPublic)

		log.Fatal(object)
	})
}
