package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/sarulabs/di"
)

const (
	DefaultOTPTTLInSeconds          = 300
	DefaultOTPLength                = 5
	DefaultAccessTokenTTLInSeconds  = 24 * 60 * 60
	DefaultRefreshTokenTTLInSeconds = 365 * 24 * 60 * 60
)

func ServiceProviderAuthentication() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AuthenticationService,
		Build: func(ctn di.Container) (interface{}, error) {
			appService := ctn.Get(service.AppService).(*app.App)
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
			otpLength := DefaultOTPLength

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

			otpLengthConfig, ok := configService.Int("authentication.otp_length")
			if ok && otpLengthConfig > 0 {
				otpLength = otpLengthConfig
			}

			passwordService := ctn.Get(service.PasswordService).(password.IPassword)
			jwtService := ctn.Get(service.JWTService).(*jwt.JWT)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			var mailService *mail.Sender
			mailServiceHitrix, err := ctn.SafeGet(service.MailMandrillService)

			if err == nil && mailServiceHitrix != nil {
				convertedMail := mailServiceHitrix.(mail.Sender)
				mailService = &convertedMail
			}

			var socialServiceMapping = make(map[string]social.IUserData)

			supportSocialLoginGoogle, ok := configService.Bool("authentication.support_social_login_google")
			if ok && supportSocialLoginGoogle {
				googleService, err := ctn.SafeGet(service.GoogleService)
				if err != nil {
					panic("google service not loaded")
				}

				socialServiceMapping[authentication.SocialLoginGoogle] = googleService.(social.IUserData)
			}

			supportSocialLoginFacebook, ok := configService.Bool("authentication.support_social_login_facebook")
			if ok && supportSocialLoginFacebook {
				googleService, err := ctn.SafeGet(service.FacebookService)
				if err != nil {
					panic("google service not loaded")
				}

				socialServiceMapping[authentication.SocialLoginFacebook] = googleService.(social.IUserData)
			}

			if appService.RedisPools == nil || appService.RedisPools.Persistent == "" {
				panic("redis persistent needs to be set")
			}

			return authentication.NewAuthenticationService(
				secret,
				accessTokenTTL,
				refreshTokenTTL,
				otpTTL,
				otpLength,
				appService,
				service.DI().Generator(),
				service.DI().ErrorLogger(),
				clockService,
				passwordService,
				jwtService,
				mailService,
				socialServiceMapping,
				service.DI().UUID(),
			), nil
		},
	}
}
