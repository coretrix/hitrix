package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/sms"
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
			configService := ctn.Get(service.ConfigService).(*config.Config)
			if configService == nil {
				panic("`authentication` key does not exists in configuration")
			}

			authConfig := configService.GetStringMap("authentication")
			if authConfig["secret"] == "" {
				panic("`authentication` key does not exists in configuration")
			}

			secret := authConfig["secret"].(string)
			accessTokenTTL := 24 * 60 * 60
			refreshTokenTTL := 365 * 24 * 60 * 60
			otpTTL := 60

			authRedis := "default"

			if authConfig["access_token_ttl"] != nil {
				accessTokenTTL = authConfig["access_token_ttl"].(int)
			}

			if authConfig["refresh_token_ttl"] != nil {
				refreshTokenTTL = authConfig["refresh_token_ttl"].(int)
			}

			if authConfig["otp_ttl"] != nil {
				otpTTL = authConfig["otp_ttl"].(int)
			}

			if authConfig["auth_redis"] != nil {
				authRedis = authConfig["auth_redis"].(string)
			}

			ormService := subContainer.Get(service.ORMEngineRequestService).(*orm.Engine)
			passwordService := ctn.Get(service.PasswordService).(*password.Password)
			jwtService := ctn.Get(service.JWTService).(*jwt.JWT)
			clockService := ctn.Get(service.ClockService).(clock.Clock)

			var smsService sms.ISender
			if authConfig["support_otp"] != nil {
				supportOTP := authConfig["support_otp"].(bool)
				if supportOTP {
					var has bool
					smsService, has = ctn.Get(service.SMSService).(sms.ISender)
					if !has {
						panic("sms service not loaded")
					}
				}
			}

			generatorService, has := service.DI().GeneratorService()
			if !has {
				panic("generator service not loaded")
			}

			return authentication.NewAuthenticationService(
				secret,
				accessTokenTTL,
				refreshTokenTTL,
				otpTTL,
				ormService,
				smsService,
				generatorService,
				clockService,
				ormService.GetRedis(authRedis),
				passwordService,
				jwtService,
			), nil
		},
	}
}
