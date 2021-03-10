package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/ddos"
	"github.com/sarulabs/di"
)

func ServiceProviderDDOS() *service.Definition {
	return &service.Definition{
		Name:   service.DDOSService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ddos.DDOS{}, nil
		},
	}
}
