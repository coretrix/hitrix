package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/sarulabs/di"
)

func FakeClockService(clockService clock.Clock) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ClockService,
		Build: func(ctn di.Container) (interface{}, error) {
			return clockService, nil
		},
	}
}
