package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/graph"
	"github.com/coretrix/hitrix/example/graph/generated"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/gin-gonic/gin"
)

func main() {
	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceDefinitionOrmRegistry(entity.Init),
		registry.ServiceDefinitionOrmEngine(),
		registry.OSSGoogle(map[string]uint64{"test": 1}),
		registry.ServiceDefinitionOrmEngineForContext(),
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
	).
		RegisterDevPanel(&entity.AdminUserEntity{}, middleware.Router, nil).Build()
	defer deferFunc()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}), func(ginEngine *gin.Engine) {
		middleware.Cors(ginEngine)
	})
}
