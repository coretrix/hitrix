package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/sarulabs/di"
)

func ServiceGoogle() *service.Definition {
	return &service.Definition{
		Name:   service.GoogleService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &social.Google{}, nil
		},
	}
}
