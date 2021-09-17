package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockSMS(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
