package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
)

func ServiceProviderMockClock(mock clock.IClock) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ClockService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
