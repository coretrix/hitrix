package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
)

func ServiceProviderClock() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ClockService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &clock.SysClock{}, nil
		},
	}
}
