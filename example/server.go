package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/graph"
	"github.com/coretrix/hitrix/example/graph/generated"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIService(
		hitrix.ServiceProviderErrorLogger(),
		hitrix.ServiceProviderConfigDirectory("config"),
		hitrix.ServiceDefinitionOrmRegistry(entity.Init),
		hitrix.ServiceDefinitionOrmEngine(),
		hitrix.ServiceDefinitionOrmEngineForContext(),
		hitrix.ServiceProviderJWT(),
		hitrix.ServiceProviderPassword(),
	).
		RegisterDevPanel(&entity.AdminUserEntity{}, middleware.Router, nil).Build()
	defer deferFunc()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}), func(ginEngine *gin.Engine) {
		middleware.Cors(ginEngine)
	})
}
