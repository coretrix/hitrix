package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
	"github.com/sarulabs/di"
)

func Crud() *service.Definition {
	return &service.Definition{
		Name:   service.CrudService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &crud.Crud{}, nil
		},
	}
}
