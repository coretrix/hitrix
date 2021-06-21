package registry

import (
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"

	"github.com/latolukasz/orm"
)

var ORMRegistryContainer []func(registry *orm.Registry)

type ORMRegistryInitFunc func(registry *orm.Registry)

func ServiceDefinitionOrmRegistry(init ORMRegistryInitFunc) *service.Definition {
	return &service.Definition{
		Name:   service.ORMConfigService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			registry := orm.NewRegistry()

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
			for _, callback := range ORMRegistryContainer {
				callback(registry)
			}

			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}
