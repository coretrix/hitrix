package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/sarulabs/di"
)

func ServiceClock() *service.Definition {
	return &service.Definition{
		Name:   service.ClockService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &clock.SysClock{}, nil
		},
	}
}
