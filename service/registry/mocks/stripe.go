package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func FakeStripeService(fake interface{}) *service.Definition {
	return &service.Definition{
		Name:   service.StripeService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
