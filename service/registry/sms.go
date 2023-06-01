package registry

import (
	"errors"

	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/sms"
)

func ServiceProviderSMS(primaryNewFunc sms.NewProviderFunc, secondaryNewFunc sms.NewProviderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.SmsTrackerEntity"]; !ok {
				return nil, errors.New("you should register SmsTrackerEntity")
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			sender := &sms.Sender{
				ConfigService:      configService,
				ClockService:       clockService,
				ErrorLoggerService: ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
			}

			var err error

			if primaryNewFunc != nil {
				sender.PrimaryProvider, err = primaryNewFunc(configService, clockService)
			}

			if err != nil {
				panic(err)
			}

			if secondaryNewFunc != nil {
				sender.SecondaryProvider, err = secondaryNewFunc(configService, clockService)
			}

			if err != nil {
				panic(err)
			}

			return sender, nil
		},
	}
}
