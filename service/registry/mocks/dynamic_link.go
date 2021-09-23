package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockDynamicLink(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.DynamicLinkService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
