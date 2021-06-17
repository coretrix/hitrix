package registry

import (
	"github.com/coretrix/hitrix/service"
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
			registry := orm.NewRegistry()

			registry.InitByYaml(service.DI().Config().Get("orm").(map[string]interface{}))
			init(registry)
			for _, callback := range ORMRegistryContainer {
				callback(registry)
			}

			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}
