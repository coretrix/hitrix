package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockUUID(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UUIDService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
