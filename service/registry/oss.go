package registry

import (
	"errors"
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/oss"
)

func ServiceProviderOSS(newFunc oss.NewProviderFunc, publicNamespaces, privateNamespaces []oss.Namespace) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)

			entities := ormConfig.GetEntities()

			if _, ok := entities["entity.OSSBucketCounterEntity"]; !ok {
				return nil, errors.New("you should register OSSBucketCounterEntity")
			}

			return newFunc(
				ctn.Get(service.ConfigService).(config.IConfig),
				ctn.Get(service.ClockService).(clock.IClock),
				publicNamespaces,
				privateNamespaces)
		},
	}
}
