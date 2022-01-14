package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockExporter(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExporterService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
