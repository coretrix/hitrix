package middleware

import (
	"github.com/coretrix/hitrix/example/rest/controller"
	"github.com/gin-gonic/gin"
)

func Router(ginEngine *gin.Engine) {
	var websocketController *controller.WebsocketController
	ginEngine.GET("/ws/", websocketController.InitConnection)
}
