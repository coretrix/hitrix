package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockCheckout(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CheckoutService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
