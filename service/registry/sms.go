package registry

import (
	"errors"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

func ServiceProviderSMS() *service.Definition {
	return &service.Definition{
		Name:   service.SMSService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			clockService := ctn.Get(service.ClockService).(clock.Clock)

			subContainer, err := ctn.SubContainer()
			if err != nil {
				return nil, err
			}

			// register twilio
			smsConfigTwilio := ctn.Get("config").(*config.Config).GetStringMap("sms.twilio")
			sid, ok := smsConfigTwilio["sid"]
			if !ok {
				return nil, errors.New("missing sms.twilio.sid")
			}

			token, ok := smsConfigTwilio["token"]
			if !ok {
				return nil, errors.New("missing sms.twilio.token")
			}

			fromNumberTwilio, ok := smsConfigTwilio["from_number"]
			if !ok {
				return nil, errors.New("missing sms.twilio.from_number")
			}

			authyURL, ok := smsConfigTwilio["authy_url"]
			if !ok {
				return nil, errors.New("missing sms.twilio.authy_url")
			}

			authyAPIKey, ok := smsConfigTwilio["authy_api_key"]
			if !ok {
				return nil, errors.New("missing sms.twilio.authy_api_key")
			}

			twilioGateway := &sms.TwilioGateway{
				SID:         sid.(string),
				Token:       token.(string),
				FromNumber:  fromNumberTwilio.(string),
				AuthyURL:    authyURL.(string),
				AuthyAPIKey: authyAPIKey.(string),
			}

			//register sinch
			smsConfigSinch := ctn.Get("config").(*config.Config).GetStringMap("sms.sinch")
			appID, ok := smsConfigSinch["app_id"]
			if !ok {
				return nil, errors.New("missing sms.sinch.app_id")
			}
			appSecret, ok := smsConfigSinch["app_secret"]
			if !ok {
				return nil, errors.New("missing sms.sinch.app_secret")
			}
			msgURL, ok := smsConfigSinch["msg_url"]
			if !ok {
				return nil, errors.New("missing sms.sinch.msg_url")
			}
			fromNumberSinch, ok := smsConfigSinch["from_number"]
			if !ok {
				return nil, errors.New("missing sms.sinch.from_number")
			}
			callURL, ok := smsConfigSinch["call_url"]
			if !ok {
				return nil, errors.New("missing sms.sinch.call_url")
			}
			callerNumber, ok := smsConfigSinch["caller_number"]
			if !ok {
				return nil, errors.New("missing sms.sinch.caller_number")
			}

			sinchGateway := &sms.SinchGateway{
				Clock:        clockService,
				AppID:        appID.(string),
				AppSecret:    appSecret.(string),
				MsgURL:       msgURL.(string),
				FromNumber:   fromNumberSinch.(string),
				CallURL:      callURL.(string),
				CallerNumber: callerNumber.(string),
			}

			// register kavenegar
			smsConfigKavenegar := ctn.Get("config").(*config.Config).GetStringMap("sms.kavenegar")
			apiKey, ok := smsConfigKavenegar["api_key"]
			if !ok {
				return nil, errors.New("missing sms.kavenegar.api_key")
			}
			sender, ok := smsConfigKavenegar["sender"]
			if !ok {
				return nil, errors.New("missing sms.kavenegar.sender")
			}

			kavenegarGateway := &sms.KavenegarGateway{
				APIKey: apiKey.(string),
				Sender: sender.(string),
			}

			ormService := subContainer.Get(service.ORMEngineRequestService).(*orm.Engine)
			registry := ormService.GetRegistry().GetSourceRegistry()

			registry.RegisterEntity(&entity.SmsTrackerEntity{})
			registry.RegisterEnumStruct("entity.SMSTrackerTypeAll", entity.SMSTrackerTypeAll)

			return &sms.Sender{
				Clock:      clockService,
				OrmService: ormService,
				GatewayFactory: map[string]sms.Gateway{
					sms.Twilio:    twilioGateway,
					sms.Sinch:     sinchGateway,
					sms.Kavenegar: kavenegarGateway,
				},
			}, nil
		},
	}
}
