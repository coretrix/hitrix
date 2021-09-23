package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/sarulabs/di"
)

func ServiceProviderMockMail(mock mail.Sender) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.MailMandrillService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
