package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockOTP(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OTPService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
