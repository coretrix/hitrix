package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

func ServiceDefinitionOrmEngine() *service.Definition {
	return serviceDefinitionOrmEngine(true)
}

func ServiceDefinitionOrmEngineForContext() *service.Definition {
	return serviceDefinitionOrmEngine(false)
}

func serviceDefinitionOrmEngine(global bool) *service.Definition {
	suffix := "request"
	if global {
		suffix = "global"
	}
	return &service.Definition{
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
