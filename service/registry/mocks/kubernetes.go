package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockKubernetes(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.KubernetesService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
