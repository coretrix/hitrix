package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockExporter(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExporterService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}