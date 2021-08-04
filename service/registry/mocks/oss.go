package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/sarulabs/di"
)

func FakeGoogleOSSService(fakeOSSService oss.Client) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSSGoogleService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fakeOSSService, nil
		},
	}
}
