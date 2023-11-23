package registry

import (
	"github.com/latolukasz/beeorm"
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
			return mail.NewSender(
				ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry),
				ctn.Get(service.ConfigService).(config.IConfig),
				ctn.Get(service.ClockService).(clock.IClock),
				ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
				newFunc,
			)
		},
	}
}
