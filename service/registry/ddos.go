package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/ddos"
	"github.com/sarulabs/di"
)

func ServiceProviderDDOS() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.DDOSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ddos.DDOS{}, nil
		},
	}
}
