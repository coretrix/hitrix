package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockGenerator(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeneratorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
