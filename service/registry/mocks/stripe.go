package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockStripe(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.StripeService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
