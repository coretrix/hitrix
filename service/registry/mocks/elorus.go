package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/elorus"
)

func ServiceProviderMockClock(mock elorus.IProvider) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ElorusService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
