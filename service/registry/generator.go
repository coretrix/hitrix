package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/sarulabs/di"
)

func GeneratorService() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeneratorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &generator.SimpleGenerator{}, nil
		},
	}
}
