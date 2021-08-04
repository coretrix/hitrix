package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeSMSService(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
