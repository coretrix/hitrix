package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/social"
)

func ServiceProviderGoogleSocial() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &social.Google{}, nil
		},
	}
}

func ServiceProviderFacebookSocial() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FacebookService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &social.Facebook{}, nil
		},
	}
}

func ServiceProviderAppleSocial() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AppleService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			return social.NewAppleSocial(configService)
		},
	}
}
