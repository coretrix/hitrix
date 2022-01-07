package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/otp"
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"
)

func ServiceProviderOTP(forceProviders ...string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OTPService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			providers := make([]otp.IOTPSMSGateway, 0)
			if len(forceProviders) > 0 {
				for _, forceProvider := range forceProviders {
					builderFunc, ok := smsOTPProviderBuilderFactory[forceProvider]
					if !ok {
						return nil, fmt.Errorf("unknown provider: %v", forceProvider)
					}

					provider, err := builderFunc(configService)
					if err != nil {
						return nil, err
					}

					providers = append(providers, provider)
				}
			} else {
				ormService := ctn.Get(service.ORMEngineGlobalService).(*beeorm.Engine)

				q := &beeorm.RedisSearchQuery{}
				q.FilterString("Key", "otp_sms_provider")

				settingsEntity := &entity.SettingsEntity{}
				if has := ormService.RedisSearchOne(settingsEntity, q); !has {
					return nil, errors.New("otp_sms_provider not found in settings")
				}

				providerKeys := strings.Split(settingsEntity.Value, ",")
				for _, providerKey := range providerKeys {
					builderFunc, ok := smsOTPProviderBuilderFactory[providerKey]
					if !ok {
						return nil, fmt.Errorf("unknown otp_sms_provider: %v", providerKey)
					}

					provider, err := builderFunc(configService)
					if err != nil {
						return nil, err
					}

					providers = append(providers, provider)
				}
			}

			if len(providers) == 0 {
				return nil, errors.New("must provide otp_sms_provider in settings or at least 1 forceProviders")
			}

			return otp.NewOTP(providers...), nil
		},
	}
}

var smsOTPProviderBuilderFactory = map[string]func(configService config.IConfig) (otp.IOTPSMSGateway, error){
	otp.SMSOTPProviderTwilio: twilioSMSOTPProviderBuilder,
	otp.SMSOTPProviderSinch:  sinchSMSOTPProviderBuilder,
}

func twilioSMSOTPProviderBuilder(configService config.IConfig) (otp.IOTPSMSGateway, error) {
	sid, ok := configService.String("sms.twilio.sid")
	if !ok {
		return nil, errors.New("missing sms.twilio.sid")
	}

	token, ok := configService.String("sms.twilio.token")

	if !ok {
		return nil, errors.New("missing sms.twilio.token")
	}

	verifySID, _ := configService.String("sms.twilio.verify_sid")

	return otp.NewTwilioSMSOTPProvider(sid, token, verifySID), nil
}

func sinchSMSOTPProviderBuilder(configService config.IConfig) (otp.IOTPSMSGateway, error) {
	appID, ok := configService.String("sms.sinch.app_id")
	if !ok {
		return nil, errors.New("missing sms.sinch.app_id")
	}
	appSecret, ok := configService.String("sms.sinch.app_secret")
	if !ok {
		return nil, errors.New("missing sms.sinch.app_secret")
	}
	verificationURL, ok := configService.String("sms.sinch.verification_url")
	if !ok {
		return nil, errors.New("missing sms.sinch.verification_url")
	}

	return otp.NewSinchSMSOTPProvider(appID, appSecret, verificationURL), nil
}
