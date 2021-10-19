package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/sarulabs/di"
)

func ServiceProviderMockOSS(mock oss.Client) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
