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

			var clientID, androidClientID string

			teamID := configService.MustString("authentication.apple.team_id")
			if clientIDInner, ok := configService.String("authentication.apple.client_id"); ok {
				clientID = clientIDInner
			}
			if androidClientIDInner, ok := configService.String("authentication.apple.android_client_id"); ok {
				androidClientID = androidClientIDInner
			}
			keyID := configService.MustString("authentication.apple.key_id")
			privateKey := configService.MustString("authentication.apple.private_key")

			return social.NewAppleSocial(
				teamID,
				clientID,
				androidClientID,
				keyID,
				privateKey,
			), nil
		},
	}
}
