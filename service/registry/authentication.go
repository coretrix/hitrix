package registry

import (
	"strconv"

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

const (
	DefaultOTPTTLInSeconds          = 60
	DefaultAccessTokenTTLInSeconds  = 24 * 60 * 60
	DefaultRefreshTokenTTLInSeconds = 365 * 24 * 60 * 60
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
				panic("`config is nil")
			}

			authConfig := configService.GetStringMapString("authentication")
			if authConfig["secret"] == "" {
				panic("`authentication` key does not exists in configuration")
			}

			secret := authConfig["secret"]
			accessTokenTTL := DefaultAccessTokenTTLInSeconds
			refreshTokenTTL := DefaultRefreshTokenTTLInSeconds
			otpTTL := DefaultOTPTTLInSeconds

			authRedis := "default"

			if authConfig["access_token_ttl"] != "" {
				accessTokenTTLInt, err := strconv.Atoi(authConfig["access_token_ttl"])
				if err == nil && accessTokenTTLInt > 0 {
					accessTokenTTL = accessTokenTTLInt
				}
			}

			if authConfig["refresh_token_ttl"] != "" {
				refreshTokenTTLInt, err := strconv.Atoi(authConfig["refresh_token_ttl"])
				if err == nil && refreshTokenTTLInt > 0 {
					refreshTokenTTL = refreshTokenTTLInt
				}
			}

			if authConfig["otp_ttl"] != "" {
				otpTTLInt, err := strconv.Atoi(authConfig["otp_ttl"])
				if err == nil && otpTTLInt > 0 {
					otpTTL = otpTTLInt
				}
			}

			if authConfig["auth_redis"] != "" {
				authRedis = authConfig["auth_redis"]
			}

			ormService := subContainer.Get(service.ORMEngineRequestService).(*orm.Engine)
			passwordService := ctn.Get(service.PasswordService).(*password.Password)
			jwtService := ctn.Get(service.JWTService).(*jwt.JWT)
			clockService := ctn.Get(service.ClockService).(clock.Clock)

			var smsService sms.ISender
			if authConfig["support_otp"] != "" {
				supportOTP := authConfig["support_otp"]
				supportOTPBool, _ := strconv.ParseBool(supportOTP)
				if supportOTPBool {
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
