package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockGoogleAnalytics(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleAnalyticsService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
