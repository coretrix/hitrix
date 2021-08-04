package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/sarulabs/di"
)

func ServiceClock() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ClockService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &clock.SysClock{}, nil
		},
	}
}
