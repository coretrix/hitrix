package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service/component/translation"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
	"github.com/coretrix/hitrix/service/component/translation"
)

func ServiceProviderCrud(exportConfigs []crud.ExportConfig) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CrudService,
		Build: func(ctn di.Container) (interface{}, error) {
			translationService, err := ctn.SafeGet(service.TranslationService)
			if err == nil {
				return &crud.Crud{ExportConfigs: exportConfigs, TranslationService: translationService.(translation.ITranslationService)}, nil
			}

			return &crud.Crud{ExportConfigs: exportConfigs}, nil
		},
	}
}
