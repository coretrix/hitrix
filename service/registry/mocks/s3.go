package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockAmazonS3(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AmazonS3Service,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
