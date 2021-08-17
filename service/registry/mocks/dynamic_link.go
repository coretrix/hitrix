package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeDynamicLinkService(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.DynamicLinkService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
