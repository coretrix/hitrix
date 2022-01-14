package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
)

func ServiceProviderCrud() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CrudService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &crud.Crud{}, nil
		},
	}
}
