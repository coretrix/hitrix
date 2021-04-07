package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/99designs/gqlgen/graphql"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/test"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func TestCreateContext(t *testing.T) {
	createContextMyApp(t, "my-app", nil, []*service.Definition{registry.APILogger(&entity.APILogEntity{})})

	apiLoggerService, has := service.DI().APILoggerService()
	if !has {
		panic("no api logger service registered")
	}

	apiLoggerService.LogStart(entity.APILogTypeApple, nil)
	apiLoggerService.LogSuccess(nil)

	apiLoggerService.LogStart(entity.APILogTypeApple, nil)
	apiLoggerService.LogError("Error appear", nil)

	ormService, _ := service.DI().OrmEngine()

	var apiLogEntities []*entity.APILogEntity
	ormService.LoadByIDs([]uint64{1, 2}, &apiLogEntities)
	assert.Len(t, apiLogEntities, 2)
	assert.Equal(t, apiLogEntities[0].Status, entity.APILogStatusCompleted)
	assert.Equal(t, apiLogEntities[1].Status, entity.APILogStatusFailed)
}

func createContextMyApp(t *testing.T, projectName string, resolvers graphql.ExecutableSchema, additionalServices []*service.Definition) *test.Environment {
	defaultServices := []*service.Definition{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceDefinitionOrmRegistry(entity.Init),
		registry.ServiceDefinitionOrmEngine(),
		registry.ServiceDefinitionOrmEngineForContext(),
	}

	return test.CreateContext(t,
		projectName,
		resolvers,
		nil,
		append(defaultServices, additionalServices...),
	)
}
