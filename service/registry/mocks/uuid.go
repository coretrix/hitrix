package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceProviderMockUUID(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UUIDService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
