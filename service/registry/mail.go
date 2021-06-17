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
			configService := ctn.Get("config").(*config.Config)

			apiKey := configService.GetString("mandrill.api_key")
			if apiKey == "" {
				panic(fmt.Errorf("mailchimp.api_key is missing"))
			}

			fromEmail := configService.GetString("mandrill.default_from_email")
			if apiKey == "" {
				panic(fmt.Errorf("mailchimp.api_key is missing"))
			}

			fromName := configService.GetString("mandrill.from_name")
			if apiKey == "" {
				panic(fmt.Errorf("mandrill.from_name is missing"))
			}

			return mail.NewMandrill(apiKey, fromEmail, fromName), nil
		},
	}
}
