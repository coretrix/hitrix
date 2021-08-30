package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	localizer "github.com/coretrix/hitrix/service/component/localizer"
	"github.com/sarulabs/di"
)

func ServiceProviderLocalizer() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LocalizerService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			apiKey, ok := configService.String("translation.poeditor.api_key")
			if !ok {
				return nil, errors.New("missing translation.poeditor.api_key")
			}
			projectId, ok := configService.String("translation.poeditor.project_id")
			if !ok {
				return nil, errors.New("missing translation.poeditor.project_id")
			}
			language, ok := configService.String("translation.poeditor.language")
			if !ok {
				return nil, errors.New("missing translation.poeditor.language")
			}

			apiSource := localizer.NewPoeditorSource(
				apiKey,
				projectId,
				language,
			)

			return localizer.NewSimpleLocalizer(apiSource), nil
		},
	}
}
