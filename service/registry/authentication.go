package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/sarulabs/di"
)

const (
	DefaultOTPTTLInSeconds          = 300
	DefaultAccessTokenTTLInSeconds  = 24 * 60 * 60
	DefaultRefreshTokenTTLInSeconds = 365 * 24 * 60 * 60
)

func ServiceProviderAuthentication() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AuthenticationService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			if configService == nil {
				panic("`config is nil")
			}

			secret, ok := configService.String("authentication.secret")
			if !ok {
				panic("secret is missing")
			}

			accessTokenTTL := DefaultAccessTokenTTLInSeconds
			refreshTokenTTL := DefaultRefreshTokenTTLInSeconds
			otpTTL := DefaultOTPTTLInSeconds

			accessTokenTTLConfig, ok := configService.Int("authentication.access_token_ttl")
			if ok && accessTokenTTLConfig > 0 {
				accessTokenTTL = accessTokenTTLConfig
			}

			refreshTokenTTLConfig, ok := configService.Int("authentication.refresh_token_ttl")
			if ok && refreshTokenTTLConfig > 0 {
				refreshTokenTTL = refreshTokenTTLConfig
			}

			otpTTLConfig, ok := configService.Int("authentication.otp_ttl")
			if ok && refreshTokenTTLConfig > 0 {
				otpTTL = otpTTLConfig
			}

			passwordService := ctn.Get(service.PasswordService).(*password.Password)
			jwtService := ctn.Get(service.JWTService).(*jwt.JWT)
			clockService := ctn.Get(service.ClockService).(clock.Clock)

			supportOTPConfig, ok := configService.Bool("authentication.support_otp")

			var smsService sms.ISender
			if ok && supportOTPConfig {
				var has bool
				smsService, has = ctn.Get(service.SMSService).(sms.ISender)
				if !has {
					panic("sms service not loaded")
				}
			}

			var mailService *mail.Sender
			mailServiceHitrix, err := ctn.SafeGet(service.MailMandrill)

			if err == nil && mailServiceHitrix != nil {
				convertedMail := mailServiceHitrix.(mail.Sender)
				mailService = &convertedMail
			}

			supportSocialLoginGoogle, ok := configService.Bool("authentication.support_social_login_google")
			var socialServiceMapping = make(map[string]social.IUserData)

			if ok && supportSocialLoginGoogle {
				googleService, err := ctn.SafeGet(service.GoogleService)
				if err != nil {
					panic("google service not loaded")
				}

				socialServiceMapping[authentication.SocialLoginGoogle] = googleService.(social.IUserData)
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
				smsService,
				generatorService,
				clockService,
				passwordService,
				jwtService,
				mailService,
				socialServiceMapping,
			), nil
		},
	}
}
