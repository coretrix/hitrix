package registry

import (
	"github.com/coretrix/hitrix/service"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
	"github.com/sarulabs/di"
)

func ServiceProviderFeatureFlag() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			return featureflag.NewFeatureFlagService(), nil
		},
	}
}
