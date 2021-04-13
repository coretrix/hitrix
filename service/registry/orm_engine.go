package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

func ServiceDefinitionOrmEngine() *service.Definition {
	return serviceDefinitionOrmEngine(true, false)
}

func ServiceDefinitionOrmEngineForContext(enableGraphQLDataLoader bool) *service.Definition {
	return serviceDefinitionOrmEngine(false, enableGraphQLDataLoader)
}

func serviceDefinitionOrmEngine(global bool, enableGraphQLDataLoader bool) *service.Definition {
	suffix := "request"
	if global {
		suffix = "global"
	}
	return &service.Definition{
		Name:   "orm_engine_" + suffix,
		Global: global,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfigService, err := ctn.SafeGet(service.ORMConfigService)
			if err != nil {
				return nil, err
			}
			ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
			if !global {
				ormEngine.EnableRequestCache(enableGraphQLDataLoader)
			}

			ormDebug := ctn.Get(service.ConfigService).(*config.Config).GetBool("orm_debug")
			if ormDebug {
				ormEngine.EnableQueryDebug()
			}

			return ormEngine, nil
		},
	}
}
