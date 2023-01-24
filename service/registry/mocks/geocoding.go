package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockGeocoding(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeocodingService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
