package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeAmazonS3Service(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AmazonS3Service,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
