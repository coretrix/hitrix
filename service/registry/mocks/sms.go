package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockSMS(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
