package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/sarulabs/di"
)

func ServiceProviderSMS(entity sms.LogEntity) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			clockService := ctn.Get(service.ClockService).(clock.Clock)
			configService := ctn.Get("config").(config.IConfig)

			// register twilio
			sid, ok := configService.String("sms.twilio.sid")
			if !ok {
				return nil, errors.New("missing sms.twilio.sid")
			}

			token, ok := configService.String("sms.twilio.token")
			if !ok {
				return nil, errors.New("missing sms.twilio.token")
			}

			fromNumberTwilio, ok := configService.String("sms.twilio.from_number")
			if !ok {
				return nil, errors.New("missing sms.twilio.from_number")
			}

			authyURL, ok := configService.String("sms.twilio.authy_url")
			if !ok {
				return nil, errors.New("missing sms.twilio.authy_url")
			}

			authyAPIKey, ok := configService.String("sms.twilio.authy_api_key")
			if !ok {
				return nil, errors.New("missing sms.twilio.authy_api_key")
			}

			verifySID, ok := configService.String("sms.twilio.verify_sid")
			if !ok {
				return nil, errors.New("missing sms.twilio.verify_sid")
			}

			twilioGateway := &sms.TwilioGateway{
				SID:         sid,
				Token:       token,
				FromNumber:  fromNumberTwilio,
				AuthyURL:    authyURL,
				AuthyAPIKey: authyAPIKey,
				VerifySID:   verifySID,
			}

			//register sinch
			appID, ok := configService.String("sms.sinch.app_id")
			if !ok {
				return nil, errors.New("missing sms.sinch.app_id")
			}
			appSecret, ok := configService.String("sms.sinch.app_secret")
			if !ok {
				return nil, errors.New("missing sms.sinch.app_secret")
			}
			msgURL, ok := configService.String("sms.sinch.msg_url")
			if !ok {
				return nil, errors.New("missing sms.sinch.msg_url")
			}
			fromNumberSinch, ok := configService.String("sms.sinch.from_number")
			if !ok {
				return nil, errors.New("missing sms.sinch.from_number")
			}
			callURL, ok := configService.String("sms.sinch.call_url")
			if !ok {
				return nil, errors.New("missing sms.sinch.call_url")
			}
			callerNumber, ok := configService.String("sms.sinch.caller_number")
			if !ok {
				return nil, errors.New("missing sms.sinch.caller_number")
			}

			sinchGateway := &sms.SinchGateway{
				Clock:        clockService,
				AppID:        appID,
				AppSecret:    appSecret,
				MsgURL:       msgURL,
				FromNumber:   fromNumberSinch,
				CallURL:      callURL,
				CallerNumber: callerNumber,
			}

			// register kavenegar
			apiKey, ok := configService.String("sms.kavenegar.api_key")
			if !ok {
				return nil, errors.New("missing sms.kavenegar.api_key")
			}
			sender, ok := configService.String("sms.kavenegar.sender")
			if !ok {
				return nil, errors.New("missing sms.kavenegar.sender")
			}

			kavenegarGateway := &sms.KavenegarGateway{
				APIKey: apiKey,
				Sender: sender,
			}

			return &sms.Sender{
				Logger: entity,
				Clock:  clockService,
				GatewayFactory: map[string]sms.Gateway{
					sms.Twilio:    twilioGateway,
					sms.Sinch:     sinchGateway,
					sms.Kavenegar: kavenegarGateway,
				},
			}, nil
		},
	}
}
