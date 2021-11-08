package registry

import (
	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
	"github.com/sarulabs/di"
)

type FeatureFlagRegistryInitFunc func(flagInterface featureflag.ServiceFeatureFlagInterface)

func ServiceProviderFeatureFlag(registry FeatureFlagRegistryInitFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			errorLoggerService := ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)
			featureFlagService := featureflag.NewFeatureFlagService(errorLoggerService)
			registry(featureFlagService)
			return featureFlagService, nil
		},
	}
}
