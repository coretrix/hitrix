package hitrix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	hitrixBinding "github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/coretrix/hitrix/service"
)

type GinInitHandler func(ginEngine *gin.Engine)
type GQLServerInitHandler func(server *handler.Server)

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

func InitGin(server graphql.ExecutableSchema, ginInitHandler GinInitHandler, gqlServerInitHandler GQLServerInitHandler) *gin.Engine {
	app := service.DI().App()
	if app.IsInProdMode() {
		gin.SetMode(gin.TestMode)
		//gin.SetMode(gin.ReleaseMode)
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

	if server != nil {
		var queryHandler gin.HandlerFunc
		if app.IsInLocalMode() || app.IsInTestMode() {
			queryHandler = graphqlHandler(server, gqlServerInitHandler)
		} else {
			timeoutSecs := service.DI().Config().DefInt64("server.timeout_sec", 10)

			queryHandler = timeout.New(
				timeout.WithTimeout(time.Duration(timeoutSecs)*time.Second),
				timeout.WithHandler(graphqlHandler(server, gqlServerInitHandler)),
				timeout.WithResponse(func(c *gin.Context) {
					service.DI().ErrorLogger().LogErrorWithRequest(c, "TIMEOUT ERROR")
				}),
			)
		}

		ginEngine.POST("/query", queryHandler)

		if app.IsInProdMode() {
			ginEngine.GET("/", middleware.AuthorizeWithQueryParam(), playgroundHandler())
		} else {
			ginEngine.GET("/", playgroundHandler())
		}
	}

	binding.Validator = hitrixBinding.NewValidator()

	return ginEngine
}

func graphqlHandler(server graphql.ExecutableSchema, gqlServerInitHandler GQLServerInitHandler) gin.HandlerFunc {
	h := handler.New(server)

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		var message string
		asErr, is := err.(error)
		if is {
			message = asErr.Error()
		} else {
			message = fmt.Sprint(err)
		}

		ginContext := ctx.Value(service.GinKey).(*gin.Context)
		requestBody := ctx.Value(service.RequestBodyKey).([]byte)
		ginContext.Request.Body = io.NopCloser(bytes.NewReader(requestBody))

		service.DI().ErrorLogger().LogErrorWithRequest(ginContext, errors.New(message))

		return errors.New("internal server error")
	})

	if gqlServerInitHandler != nil {
		gqlServerInitHandler(h)
	}

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		c.Writer.Header().Add("X-Robots-Tag", "noindex")
		h.ServeHTTP(c.Writer, c.Request)
	}
}
