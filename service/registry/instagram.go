package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/instagram"
)

func ServiceProviderInstagram(newFunctions instagram.NewProviderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.InstagramService,
		Build: func(ctn di.Container) (interface{}, error) {
			return instagram.NewAPIManager(
				ctn.Get(service.ConfigService).(config.IConfig),
				newFunctions)
		},
	}
}
