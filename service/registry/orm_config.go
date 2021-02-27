package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"

	"github.com/summer-solutions/orm"
)

type ORMRegistryInitFunc func(registry *orm.Registry)

func ServiceDefinitionOrmRegistry(init ORMRegistryInitFunc) *service.Definition {
	return &service.Definition{
		Name:   "orm_config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			registry := orm.InitByYaml(service.DI().Config().Get("orm").(map[string]interface{}))
			init(registry)
			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}