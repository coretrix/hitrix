package registry

import (
	"context"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/latolukasz/beeorm"
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

			ormEngine := ormConfigService.(beeorm.ValidatedRegistry).CreateEngine(context.Background())
			if !global && enableGraphQLDataLoader {
				ormEngine.EnableRequestCache()
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)

			ormDebug, ok := configService.Bool("orm_debug")
			if ok && ormDebug {
				ormEngine.EnableQueryDebug()
			}

			return ormEngine, nil
		},
	}
}
