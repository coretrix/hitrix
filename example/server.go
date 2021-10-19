package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/graph"
	"github.com/coretrix/hitrix/example/graph/generated"
	model "github.com/coretrix/hitrix/example/model/socket"
	exampleMiddleware "github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/gin-gonic/gin"
)

var eventHandlersMap = socket.NamespaceEventHandlerMap{
	model.DefaultWebsocketNamespace: &socket.EventHandlers{RegisterHandler: model.RegisterSocketHandler, UnregisterHandler: model.UnRegisterSocketHandler},
}

func main() {
	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrmEngine(),
		registry.ServiceProviderOSS(map[string]uint64{"test": 1}),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
		registry.ServiceProviderSocketRegistry(eventHandlersMap),
		registry.ServiceProviderOTP(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false),
	).RegisterRedisPools(&app.RedisPools{Persistent: "default", Cache: "default"}).
		RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.Router).Build()
	defer deferFunc()

	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}), func(ginEngine *gin.Engine) {
		exampleMiddleware.Router(ginEngine)
		middleware.Cors(ginEngine)
	}, nil)
}
