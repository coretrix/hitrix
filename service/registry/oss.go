package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service/component/clock"

	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"
)

func ServiceProviderOSS(bucketsMapping map[string]*oss.Bucket, ossProvider int) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSService,
		Build: func(ctn di.Container) (interface{}, error) {
			if len(bucketsMapping) == 0 {
				panic("buckets mapping not defined")
			}

			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)

			entities := ormConfig.GetEntities()

			if _, ok := entities["entity.OSSBucketCounterEntity"]; !ok {
				return nil, errors.New("you should register OSSBucketCounterEntity")
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			if ossProvider == oss.ProviderGoogleOSS {
				return oss.NewGoogleOSS(
					configService,
					clockService,
					bucketsMapping,
					ctn.Get(service.AppService).(*app.App).Mode)
			} else if ossProvider == oss.ProviderAmazonOSS {
				return oss.NewAmazonOSS(
					configService,
					clockService,
					bucketsMapping,
					ctn.Get(service.AppService).(*app.App).Mode)
			}

			panic("invalid oss provider")
		},
	}
}
