package registry

import (
	"errors"

	"github.com/latolukasz/beeorm/v2"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/mail"
)

// ServiceProviderMail Be sure that you registered entity MailTrackerEntity
func ServiceProviderMail(newFunc mail.NewSenderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.MailService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.MailTrackerEntity"]; !ok {
				return nil, errors.New("you should register MailTrackerEntity")
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			provider, err := newFunc(ctn.Get(service.ConfigService).(config.IConfig))
			if err != nil {
				return nil, err
			}

			return &mail.Sender{
				ConfigService:      configService,
				ClockService:       clockService,
				Provider:           provider,
				ErrorLoggerService: ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
			}, nil
		},
	}
}
