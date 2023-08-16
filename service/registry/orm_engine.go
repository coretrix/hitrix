package registry

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm/v2"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
)

func ServiceProviderOrmEngine(searchPool ...string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ORMEngineGlobalService,
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

			dataLayer := &datalayer.DataLayer{
				Engine: ormEngine,
			}

			if len(searchPool) != 0 && searchPool[0] != "" {
				dataLayer.RedisSearch = redisearch.NewRedisSearch(service.DI().App().GlobalContext, ormEngine, searchPool[0])
			}

			return dataLayer, nil
		},
	}
}

func ServiceProviderOrmEngineForContext(enableGraphQLDataLoader bool, searchPool ...string) *service.DefinitionRequest {
	return &service.DefinitionRequest{
		Name: service.ORMEngineRequestService,
		Build: func(c *gin.Context) (interface{}, error) {
			ormConfigService := service.DI().OrmConfig()

			ormEngine := ormConfigService.CreateEngine()
			if enableGraphQLDataLoader {
				ormEngine.EnableRequestCache()
			}

			configService := service.DI().Config()

			ormDebug, ok := configService.Bool("orm_debug")
			if ok && ormDebug {
				ormEngine.EnableQueryDebug()
			}

			dataLayer := &datalayer.DataLayer{
				Engine: ormEngine,
			}

			if len(searchPool) != 0 && searchPool[0] != "" {
				dataLayer.RedisSearch = redisearch.NewRedisSearch(c, ormEngine, searchPool[0])
			}

			return dataLayer, nil
		},
	}
}
