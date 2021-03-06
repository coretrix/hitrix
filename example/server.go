package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/graph"
	"github.com/coretrix/hitrix/example/graph/generated"
	model "github.com/coretrix/hitrix/example/model/socket"
	exampleMiddleware "github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/middleware"
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
	).RegisterDIService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceDefinitionOrmRegistry(entity.Init),
		registry.ServiceDefinitionOrmEngine(),
		registry.OSSGoogle(map[string]uint64{"test": 1}),
		registry.ServiceDefinitionOrmEngineForContext(false),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
		registry.ServiceSocketRegistry(eventHandlersMap),
	).
		RegisterDevPanel(&entity.AdminUserEntity{}, middleware.Router, nil, nil).Build()
	defer deferFunc()

	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}), func(ginEngine *gin.Engine) {
		exampleMiddleware.Router(ginEngine)
		middleware.Cors(ginEngine)
	}, nil)
}
