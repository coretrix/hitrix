package registry

import (
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"

	"github.com/sarulabs/di"
)

// MailMandrill Be sure that you registered entity MailTrackerEntity
func MailMandrill() *service.Definition {
	return &service.Definition{
		Name:   service.MailMandrill,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.MailTrackerEntity"]; !ok {
				return nil, errors.New("you should register MailTrackerEntity")
			}

			configService := ctn.Get("config").(config.IConfig)

			apiKey, ok := configService.String("mandrill.api_key")
			if !ok {
				panic(fmt.Errorf("mailchimp.api_key is missing"))
			}

			fromEmail, ok := configService.String("mandrill.default_from_email")
			if !ok {
				panic(fmt.Errorf("mailchimp.default_from_email is missing"))
			}

			fromName, ok := configService.String("mandrill.from_name")
			if !ok {
				panic(fmt.Errorf("mandrill.from_name is missing"))
			}

			return mail.NewMandrill(apiKey, fromEmail, fromName), nil
		},
	}
}
