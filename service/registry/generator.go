package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/sarulabs/di"
)

func GeneratorService() *service.Definition {
	return &service.Definition{
		Name:   service.GeneratorService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &generator.SimpleGenerator{}, nil
		},
	}
}
