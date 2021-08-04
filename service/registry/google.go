package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/sarulabs/di"
)

func ServiceGoogle() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &social.Google{}, nil
		},
	}
}
