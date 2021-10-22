package hitrix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/timeout"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GinInitHandler func(ginEngine *gin.Engine)
type GQLServerInitHandler func(server *handler.Server)

func contextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			panic(err)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

		ctx := context.WithValue(c.Request.Context(), service.GinKey, c)
		ctx = context.WithValue(ctx, service.RequestBodyKey, body)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func InitGin(server graphql.ExecutableSchema, ginInitHandler GinInitHandler, gqlServerInitHandler GQLServerInitHandler) *gin.Engine {
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

	if server != nil {
		var queryHandler gin.HandlerFunc
		if app.IsInLocalMode() {
			queryHandler = graphqlHandler(server, gqlServerInitHandler)
		} else {
			queryHandler = timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server, gqlServerInitHandler)))
		}

		ginEngine.POST("/query", queryHandler)
		ginEngine.GET("/", playgroundHandler())
	}
	binding.Validator = helper.NewValidator()
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
		errorMessage := message + "\n" + string(debug.Stack())

		ginContext := ctx.Value(service.GinKey).(*gin.Context)
		requestBody := ctx.Value(service.RequestBodyKey).([]byte)
		ginContext.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

		service.DI().ErrorLogger().LogErrorWithRequest(ginContext, errors.New(errorMessage))

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
