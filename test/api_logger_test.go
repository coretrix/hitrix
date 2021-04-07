package main

import (
	"testing"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/stretchr/testify/assert"
)

func TestApiLogger(t *testing.T) {
	createContextMyApp(t, "my-app", nil, registry.APILogger(&entity.APILogEntity{}))

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
