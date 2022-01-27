package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/template"
)

func ServiceProviderTemplate() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.TemplateService,
		Build: func(ctn di.Container) (interface{}, error) {
			return template.NewTemplateService(), nil
		},
	}
}
