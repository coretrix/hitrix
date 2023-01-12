package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockLicensePlateRecognizer(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LicensePlateRecognizerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
