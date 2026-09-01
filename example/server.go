package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	model "github.com/coretrix/hitrix/example/model/socket"
	exampleOSS "github.com/coretrix/hitrix/example/oss"
	exampleMiddleware "github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/oss"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/coretrix/hitrix/service/registry"
)

var eventHandlersMap = socket.NamespaceEventHandlerMap{
	model.DefaultWebsocketNamespace: &socket.EventHandlers{
		RegisterHandler:   model.RegisterSocketHandler,
		UnregisterHandler: model.UnRegisterSocketHandler,
	},
}

func main() {
	socketConfig := socket.DefaultConfig()
	socketConfig.CheckOrigin = func(request *http.Request) bool {
		origin := request.Header.Get("Origin")

		return origin == "" || origin == "http://localhost:9999" || origin == "http://localhost:63342"
	}

	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrmEngine(),
		registry.ServiceProviderClock(),
		registry.ServiceProviderOSS(oss.NewAmazonOSS, exampleOSS.Namespaces),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(password.NewSimpleManager),
		registry.ServiceProviderSocketRegistryWithConfig(eventHandlersMap, socketConfig),
		registry.ServiceProviderOTP(nil),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(),
	).RegisterRedisPools(
		&app.RedisPools{
			Persistent: "default",
			Cache:      "default",
			Search:     []string{"search_pool", "search_pool2"},
			Stream:     "stream_pool",
		},
	).RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.DevPanelRouter).Build()
	defer deferFunc()

	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()
	b.RunAsyncRequestLoggerCleaner()

	s.RunServer(9999, func(ginEngine *gin.Engine) {
		middleware.RequestLogger(ginEngine, nil)
		exampleMiddleware.Router(ginEngine)
		middleware.Cors(ginEngine)
	})
}
