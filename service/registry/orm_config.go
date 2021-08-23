package registry

import (
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"

	"github.com/latolukasz/beeorm"
)

type ORMRegistryInitFunc func(registry *beeorm.Registry)

func ServiceDefinitionOrmRegistry(init ORMRegistryInitFunc) *service.DefinitionGlobal {
	var defferFunc func()
	var ormConfig beeorm.ValidatedRegistry
	var err error

	return &service.DefinitionGlobal{
		Name: service.ORMConfigService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

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

			ormConfig, defferFunc, err = registry.Validate()
			return ormConfig, err
		},
		Close: func(obj interface{}) error {
			defferFunc()
			return nil
		},
	}
}
