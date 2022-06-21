package registry

import (
	"errors"

	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
)

type FeatureFlagRegistryInitFunc func(flagInterface featureflag.ServiceFeatureFlagInterface)

func ServiceProviderFeatureFlag(registry FeatureFlagRegistryInitFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.FeatureFlagEntity"]; !ok {
				return nil, errors.New("you should register FeatureFlagEntity")
			}

			errorLoggerService := ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)
			featureFlagService := featureflag.NewFeatureFlagService(errorLoggerService)
			registry(featureFlagService)
			return featureFlagService, nil
		},
	}
}
func ServiceProviderFeatureFlagWithCache(registry FeatureFlagRegistryInitFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.FeatureFlagEntity"]; !ok {
				return nil, errors.New("you should register FeatureFlagEntity")
			}

			errorLoggerService := ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			featureFlagService := featureflag.NewFeatureFlagWithCacheService(errorLoggerService, clockService)
			registry(featureFlagService)
			return featureFlagService, nil
		},
	}
}
