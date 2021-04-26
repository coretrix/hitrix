package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

func ServiceProviderAuthentication() *service.Definition {
	return &service.Definition{
		Name:   service.AuthenticationService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			subContainer, err := ctn.SubContainer()
			if err != nil {
				return nil, err
			}
			config := service.DI().Config().Get("authentication")
			if config == nil {
				panic("`authentication` key does not exists in configuration")
			}
			configMap := config.(map[string]interface{})
			if configMap["secret"] == nil {
				panic("`authentication` key does not exists in configuration")
			}

			secret := configMap["secret"].(string)
			accessTokenTTL := 24 * 60 * 60
			refreshTokenTTL := 365 * 24 * 60 * 60

			authRedis := "default"

			if configMap["accessTokenTTL"] != nil {
				accessTokenTTL = configMap["accessTokenTTL"].(int)
			}

			if configMap["refreshTokenTTL"] != nil {
				refreshTokenTTL = configMap["refreshTokenTTL"].(int)
			}

			if configMap["authRedis"] != nil {
				authRedis = configMap["authRedis"].(string)
			}

			ormService := subContainer.Get(service.ORMEngineRequestService).(*orm.Engine)
			passwordService := subContainer.Get(service.PasswordService).(*password.Password)
			jwtService := subContainer.Get(service.JWTService).(*jwt.JWT)
			return authentication.NewAuthenticationService(
				secret,
				accessTokenTTL,
				refreshTokenTTL,
				ormService,
				ormService.GetRedis(authRedis),
				passwordService,
				jwtService,
			), nil
		},
	}
}
