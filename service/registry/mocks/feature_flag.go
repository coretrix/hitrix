package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceFeatureFlag(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
