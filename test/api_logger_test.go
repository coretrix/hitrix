package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func TestApiLogger(t *testing.T) {
	createContextMyApp(t, "server", nil,
		[]*service.DefinitionGlobal{
			registry.ServiceProviderAPILogger(&entity.APILogEntity{}),
		},
		nil,
	)

	apiLoggerService := service.DI().APILogger()

	ormService := service.DI().OrmEngine()
	apiLoggerService.LogStart(ormService, entity.APILogTypeApple, nil)
	apiLoggerService.LogSuccess(ormService, nil)

	apiLoggerService.LogStart(ormService, entity.APILogTypeApple, nil)
	apiLoggerService.LogError(ormService, "Error appear", nil)

	var apiLogEntities []*entity.APILogEntity
	ormService.LoadByIDs([]uint64{1, 2}, &apiLogEntities)
	assert.Len(t, apiLogEntities, 2)
	assert.Equal(t, apiLogEntities[0].Status, entity.APILogStatusCompleted)
	assert.Equal(t, apiLogEntities[1].Status, entity.APILogStatusFailed)
}
