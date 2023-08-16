package main

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"

	"github.com/coretrix/hitrix/example/entity/initialize"
	"github.com/coretrix/hitrix/example/redis"
	"github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/registry"
)

func createContextMyApp(
	t *testing.T,
	projectName string, //nolint //`projectName` always receives `"my-app"`
	resolvers graphql.ExecutableSchema, //nolint //`resolvers` always receives `nil`
	mockGlobalServices []*service.DefinitionGlobal,
	mockRequestServices []*service.DefinitionRequest, //nolint //`resolvers` always receives `nil`
) *test.Environment {
	defaultGlobalServices := []*service.DefinitionGlobal{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceProviderOrmRegistry(initialize.Init),
		registry.ServiceProviderCrud(nil),
		registry.ServiceProviderOrmEngine(redis.SearchPool),
		registry.ServiceProviderErrorLogger(),
	}

	defaultRequestServices := []*service.DefinitionRequest{
		registry.ServiceProviderOrmEngineForContext(false, redis.SearchPool),
	}

	return test.CreateAPIContext(t,
		projectName,
		resolvers,
		middleware.Router,
		defaultGlobalServices,
		defaultRequestServices,
		mockGlobalServices,
		mockRequestServices,
		&app.RedisPools{
			Cache:      redis.DefaultPool,
			Persistent: redis.DefaultPool,
			Stream:     redis.StreamsPool,
			Search:     redis.SearchPool,
		},
	)
}
