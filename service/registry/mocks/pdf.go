package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeServiceTemplate(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.PDFService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
