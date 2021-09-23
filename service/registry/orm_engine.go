package registry

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"
)

func ServiceProviderOrmEngine() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: "orm_engine_global",

		Build: func(ctn di.Container) (interface{}, error) {
			ormConfigService, err := ctn.SafeGet(service.ORMConfigService)
			if err != nil {
				return nil, err
			}

			ormEngine := ormConfigService.(beeorm.ValidatedRegistry).CreateEngine()

			configService := ctn.Get(service.ConfigService).(config.IConfig)

			ormDebug, ok := configService.Bool("orm_debug")
			if ok && ormDebug {
				ormEngine.EnableQueryDebug()
			}

			return ormEngine, nil
		},
	}
}

func ServiceProviderOrmEngineForContext(enableGraphQLDataLoader bool) *service.DefinitionRequest {
	return &service.DefinitionRequest{
		Name: "orm_engine_request",
		Build: func(c *gin.Context) (interface{}, error) {
			ormConfigService, has := service.DI().OrmConfig()
			if !has {
				return nil, errors.New("not registered orm config")
			}

			ormEngine := ormConfigService.CreateEngine()
			if enableGraphQLDataLoader {
				ormEngine.EnableRequestCache()
			}

			configService := service.DI().Config()

			ormDebug, ok := configService.Bool("orm_debug")
			if ok && ormDebug {
				ormEngine.EnableQueryDebug()
			}

			return ormEngine, nil
		},
	}
}
