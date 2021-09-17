package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/uuid"
	"github.com/sarulabs/di"
)

func ServiceProviderUUID() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UUIDService,
		Build: func(ctn di.Container) (interface{}, error) {
			return uuid.NewGoogleUUID(), nil
		},
	}
}
