package main

import (
	"testing"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/rest/middleware"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/registry"
)

func createContextMyApp(
	t *testing.T,
	projectName string, //nolint //`projectName` always receives `"my-app"`
	mockGlobalServices []*service.DefinitionGlobal,
	mockRequestServices []*service.DefinitionRequest, //nolint //`resolvers` always receives `nil`
) *test.Environment {
	defaultGlobalServices := []*service.DefinitionGlobal{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderCrud(nil),
		registry.ServiceProviderOrmEngine(),
	}

	defaultRequestServices := []*service.DefinitionRequest{
		registry.ServiceProviderOrmEngineForContext(),
	}

	return test.CreateAPIContext(t,
		projectName,
		middleware.Router,
		defaultGlobalServices,
		defaultRequestServices,
		mockGlobalServices,
		mockRequestServices,
		&app.RedisPools{Persistent: "default", Cache: "default", Search: []string{"search_pool", "search_pool2"}, Stream: "stream_pool"},
	)
}
