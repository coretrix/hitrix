package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"
)

func ServiceProviderConfigDirectory(configDirectory string) *service.Definition {
	return &service.Definition{
		Name:   "config_directory",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return configDirectory, nil
		},
	}
}

func ServiceConfig() *service.Definition {
	return &service.Definition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configDirectory := ctn.Get("config_directory").(string)
			return config.NewViperConfig(service.DI().App().Name, service.DI().App().Mode, configDirectory)
		},
	}
}
