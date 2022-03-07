package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockInstagram(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.InstagramService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
