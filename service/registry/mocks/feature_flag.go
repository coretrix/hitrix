package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceFeatureFlag(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
