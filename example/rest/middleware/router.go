package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/example/rest/controller"
)

func Router(ginEngine *gin.Engine) {
	var websocketController *controller.WebsocketController
	ginEngine.GET("/ws/", websocketController.InitConnection)
}
