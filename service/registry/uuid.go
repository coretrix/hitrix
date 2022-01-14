package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/uuid"
)

func ServiceProviderUUID() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UUIDService,
		Build: func(ctn di.Container) (interface{}, error) {
			return uuid.NewGoogleUUID(), nil
		},
	}
}
