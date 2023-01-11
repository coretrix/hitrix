package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockLicencePlateRecognizer(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LicencePlateRecognizerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
