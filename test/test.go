package main

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func createContextMyApp(t *testing.T, projectName string, resolvers graphql.ExecutableSchema, mockGlobalServices []*service.DefinitionGlobal, mockRequestServices []*service.DefinitionRequest) *test.Environment {
	defaultGlobalServices := []*service.DefinitionGlobal{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceDefinitionOrmRegistry(entity.Init),
		registry.ServiceDefinitionOrmEngine(),
	}

	defaultRequestServices := []*service.DefinitionRequest{
		registry.ServiceDefinitionOrmEngineForContext(false),
	}

	return test.CreateAPIContext(t,
		projectName,
		resolvers,
		nil,
		defaultGlobalServices,
		defaultRequestServices,
		mockGlobalServices,
		mockRequestServices,
	)
}
