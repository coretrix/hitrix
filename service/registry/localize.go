package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/localize"
	"github.com/sarulabs/di"
)

func ServiceProviderLocalize() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LocalizeService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			apiKey, ok := configService.String("translation.poeditor.api_key")
			if !ok {
				return nil, errors.New("missing translation.poeditor.api_key")
			}
			projectID, ok := configService.String("translation.poeditor.project_id")
			if !ok {
				return nil, errors.New("missing translation.poeditor.project_id")
			}
			language, ok := configService.String("translation.poeditor.language")
			if !ok {
				return nil, errors.New("missing translation.poeditor.language")
			}

			apiSource := localize.NewPoeditorSource(
				apiKey,
				projectID,
				language,
			)

			return localize.NewSimpleLocalizer(apiSource), nil
		},
	}
}
