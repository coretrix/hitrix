package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeAmazonS3Service(fake interface{}) *service.Definition {
	return &service.Definition{
		Name:   service.AmazonS3Service,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
