package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeSMSService(fake interface{}) *service.Definition {
	return &service.Definition{
		Name:   service.SMSService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
