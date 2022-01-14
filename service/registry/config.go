package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
)

func ServiceProviderConfigDirectory(configDirectory string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: "config_directory",
		Build: func(ctn di.Container) (interface{}, error) {
			return configDirectory, nil
		},
	}
}

func ServiceProviderConfig() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ConfigService,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewConfig(service.DI().App().Name, service.DI().App().Mode, ctn.Get("config_directory").(string))
		},
	}
}
