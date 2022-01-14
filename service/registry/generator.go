package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/generator"
)

func ServiceProviderGenerator() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeneratorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &generator.SimpleGenerator{}, nil
		},
	}
}
