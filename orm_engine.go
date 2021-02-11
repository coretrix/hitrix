package hitrix

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

func ServiceDefinitionOrmEngine() *ServiceDefinition {
	return serviceDefinitionOrmEngine(true)
}

func ServiceDefinitionOrmEngineForContext() *ServiceDefinition {
	return serviceDefinitionOrmEngine(false)
}

func serviceDefinitionOrmEngine(global bool) *ServiceDefinition {
	suffix := "request"
	if global {
		suffix = "global"
	}
	return &ServiceDefinition{
		Name:   "orm_engine_" + suffix,
		Global: global,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfigService, err := ctn.SafeGet("orm_config")
			if err != nil {
				return nil, err
			}
			ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
			if !global {
				ormEngine.EnableRequestCache(true)
			}
			return ormEngine, nil
		},
	}
}
