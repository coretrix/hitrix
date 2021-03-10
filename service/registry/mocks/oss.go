package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/sarulabs/di"
)

func FakeGoogleOSSService(fakeOSSService oss.Client) *service.Definition {
	return &service.Definition{
		Name:   service.OSSGoogleService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &fakeOSSService, nil
		},
	}
}
