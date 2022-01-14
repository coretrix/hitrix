package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockCheckout(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CheckoutService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
