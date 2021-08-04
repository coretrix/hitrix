package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/sarulabs/di"
)

func FakeMailService(fakeMailSender mail.Sender) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.MailMandrill,
		Build: func(ctn di.Container) (interface{}, error) {
			return fakeMailSender, nil
		},
	}
}
