package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockGenerator(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeneratorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
