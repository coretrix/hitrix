package registry

import (
	"errors"

	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
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

			return newFunc(ctn.Get(service.ConfigService).(config.IConfig))
		},
	}
}
