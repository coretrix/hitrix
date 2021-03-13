package main

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func TestCreateContext(t *testing.T) {
	createContextWebAPI(t, "my-app", nil)
}

func createContextWebAPI(t *testing.T, projectName string, resolvers graphql.ExecutableSchema) *test.Ctx {
	defaultServices := []*service.Definition{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceDefinitionOrmRegistry(entity.Init),
		registry.ServiceDefinitionOrmEngine(),
	}

	return test.CreateContext(t,
		projectName,
		resolvers,
		defaultServices,
	)
}
