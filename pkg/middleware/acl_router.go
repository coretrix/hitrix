package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/controller"
)

func FileRouter(ginEngine *gin.Engine) {
	v1Group := ginEngine.Group("/v1/")

	var fileController *controller.FileController
	fileGroup := v1Group.Group("file/")
	{
		fileGroup.POST("upload/", fileController.PostUploadImageAction)
	}
}
