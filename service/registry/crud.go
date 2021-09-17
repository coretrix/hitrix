package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
	"github.com/sarulabs/di"
)

func ServiceProviderCrud() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CrudService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &crud.Crud{}, nil
		},
	}
}
