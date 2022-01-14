package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"
)

func ServiceProviderMockOSS(mock oss.IProvider) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
