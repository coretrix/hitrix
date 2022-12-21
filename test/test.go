package main

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/registry"
)

func createContextMyApp(
	t *testing.T,
	//nolint //`projectName` always receives `"my-app"`
	projectName string,
	//nolint //`resolvers` always receives `nil`
	resolvers graphql.ExecutableSchema,
	mockGlobalServices []*service.DefinitionGlobal,
	//nolint //`mockRequestServices` always receives `nil`
	mockRequestServices []*service.DefinitionRequest,
) *test.Environment {
	defaultGlobalServices := []*service.DefinitionGlobal{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderCrud(),
		registry.ServiceProviderOrmEngine(),
	}

	defaultRequestServices := []*service.DefinitionRequest{
		registry.ServiceProviderOrmEngineForContext(false),
	}

	return test.CreateAPIContext(t,
		projectName,
		resolvers,
		middleware.Router,
		defaultGlobalServices,
		defaultRequestServices,
		mockGlobalServices,
		mockRequestServices,
		&app.RedisPools{Cache: "default", Persistent: "default"},
	)
}
