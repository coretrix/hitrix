package main

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/entity/initialize"
	"github.com/coretrix/hitrix/example/graph"
	"github.com/coretrix/hitrix/example/graph/generated"
	model "github.com/coretrix/hitrix/example/model/socket"
	exampleOSS "github.com/coretrix/hitrix/example/oss"
	"github.com/coretrix/hitrix/example/redis"
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
	s, deferFunc := hitrix.New(
		"server", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("../../config"),
		registry.ServiceProviderOrmRegistry(initialize.Init),
		registry.ServiceProviderOrmEngine(redis.SearchPool),
		registry.ServiceProviderClock(),
		registry.ServiceProviderOSS(oss.NewAmazonOSS, exampleOSS.Namespaces),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(password.NewSimpleManager),
		registry.ServiceProviderSocketRegistry(eventHandlersMap),
		registry.ServiceProviderOTP(),
		registry.ServiceProviderRequestLogger(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false, redis.SearchPool),
	).RegisterRedisPools(&app.RedisPools{
		Cache:      redis.DefaultPool,
		Persistent: redis.DefaultPool,
		Stream:     redis.StreamsPool,
		Search:     redis.SearchPool,
	}).RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.DevPanelRouter).Build()
	defer deferFunc()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}), func(ginEngine *gin.Engine) {
		//middleware.RequestLogger(ginEngine, nil)
		exampleMiddleware.Router(ginEngine)
		middleware.Cors(ginEngine)
	}, nil)
}
