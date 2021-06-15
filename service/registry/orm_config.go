package registry

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

type ORMRegistryInitFunc func(registry *orm.Registry)

func ServiceDefinitionOrmRegistry(init ORMRegistryInitFunc, customRedisSearchIndexes ...*orm.RedisSearchIndex) *service.Definition {
	return &service.Definition{
		Name:   service.ORMConfigService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			registry := orm.NewRegistry()

			registry.InitByYaml(service.DI().Config().Get("orm").(map[string]interface{}))
			init(registry)

			_, err := ctn.SafeGet(service.OSSGoogleService)
			if err == nil {
				registry.RegisterEntity(&entity.OSSBucketCounterEntity{})
			}

			_, err = ctn.SafeGet(service.AmazonS3Service)
			if err == nil {
				registry.RegisterEntity(&entity.S3BucketCounterEntity{})
			}

			_, err = ctn.SafeGet(service.MailMandrill)
			if err == nil {
				registry.RegisterEntity(&entity.MailTrackerEntity{})
				registry.RegisterEnumStruct("entity.MailTrackerStatusAll", entity.MailTrackerStatusAll)
			}

			for _, customRedisSearchIndex := range customRedisSearchIndexes {
				registry.RegisterRedisSearchIndex(customRedisSearchIndex)
			}

			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}
