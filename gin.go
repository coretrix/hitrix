package hitrix

import (
	"bytes"
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	hitrixBinding "github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/service"
)

type GinInitHandler func(ginEngine *gin.Engine)

func contextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			panic(err)
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		ctx := context.WithValue(c.Request.Context(), service.GinKey, c)
		ctx = context.WithValue(ctx, service.RequestBodyKey, body)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func InitGin(ginInitHandler GinInitHandler) *gin.Engine {
	app := service.DI().App()
	if app.IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	} else if app.IsInTestMode() {
		gin.SetMode(gin.TestMode)
	}

	ginEngine := gin.New()

	if !app.IsInProdMode() {
		ginEngine.Use(gin.Logger())
	}

	ginEngine.Use(recovery())
	ginEngine.Use(contextToContextMiddleware())

	if ginInitHandler != nil {
		ginInitHandler(ginEngine)
	}

	if app.DevPanel != nil {
		devRouter := app.DevPanel.Router
		devRouter(ginEngine)
	}

	binding.Validator = hitrixBinding.NewValidator()

	return ginEngine
}
