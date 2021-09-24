package registry

import (
	"errors"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/otp"
	"github.com/sarulabs/di"
)

func ServiceProviderOTP() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OTPService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get("config").(config.IConfig)

			sid, ok := configService.String("sms.twilio.sid")
			if !ok {
				return nil, errors.New("missing sms.twilio.sid")
			}

			token, ok := configService.String("sms.twilio.token")

			if !ok {
				return nil, errors.New("missing sms.twilio.token")
			}

			verifySID, _ := configService.String("sms.twilio.verify_sid")

			return otp.NewOTP(otp.NewTwilioSMSOTPProvider(sid, token, verifySID)), nil
		},
	}
}
