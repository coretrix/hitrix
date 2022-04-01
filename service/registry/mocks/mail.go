package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/mail"
)

func ServiceProviderMockMail(mock mail.Sender) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.MailService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
