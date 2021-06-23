package registry

import (
	"fmt"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/latolukasz/orm"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"

	"github.com/sarulabs/di"
)

func MailMandrill() *service.Definition {
	ORMRegistryContainer = append(ORMRegistryContainer, func(registry *orm.Registry) {
		registry.RegisterEntity(&entity.MailTrackerEntity{})
		registry.RegisterEnumStruct("entity.MailTrackerStatusAll", entity.MailTrackerStatusAll)
	})

	return &service.Definition{
		Name:   service.MailMandrill,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
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
