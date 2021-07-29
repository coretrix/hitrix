package registry

import (
	"errors"
	"fmt"
	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"

	"github.com/latolukasz/beeorm"
)

type ORMRegistryInitFunc func(registry *beeorm.Registry)

func ServiceDefinitionOrmRegistry(init ORMRegistryInitFunc) *service.Definition {
	return &service.Definition{
		Name:   service.ORMConfigService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			appService := ctn.Get(service.AppService).(*app.App)

			registry := beeorm.NewRegistry()

			configuration, ok := configService.Get("orm")
			if !ok {
				return nil, errors.New("no orm config")
			}

			yamlConfig := map[string]interface{}{}
			for k, v := range configuration.(map[interface{}]interface{}) {
				yamlConfig[fmt.Sprint(k)] = v
			}

			registry.InitByYaml(yamlConfig)
			init(registry)

			ormConfig, err := registry.Validate(appService.GlobalContext)
			return ormConfig, err
		},
	}
}
