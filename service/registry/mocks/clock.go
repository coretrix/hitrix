package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/sarulabs/di"
)

func FakeClockService(clockService clock.Clock) *service.Definition {
	return &service.Definition{
		Name:   service.ClockService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return clockService, nil
		},
	}
}
