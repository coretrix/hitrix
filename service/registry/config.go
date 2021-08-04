package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"
)

func ServiceProviderConfigDirectory(configDirectory string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: "config_directory",
		Build: func(ctn di.Container) (interface{}, error) {
			return configDirectory, nil
		},
	}
}

func ServiceConfig() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ConfigService,

		Build: func(ctn di.Container) (interface{}, error) {
			configDirectory := ctn.Get("config_directory").(string)
			return config.NewConfig(service.DI().App().Name, service.DI().App().Mode, configDirectory)
		},
	}
}
