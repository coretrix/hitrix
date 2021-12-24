package registry

import (
	"errors"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"log"
	"os"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/localize"
	"github.com/sarulabs/di"
)

func ServiceProviderLocalize(projectNameEnvVar string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LocalizeService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			var apiSource localize.Source

			if _, ok := configService.StringMap("translation.poeditor"); ok {
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

				apiSource = localize.NewPoeditorSource(
					apiKey,
					projectID,
					language,
				)
			}

			var path string
			if projectNameEnvVar != "" {
				path = configService.GetFolderPath() + "/../locale/" + os.Getenv(projectNameEnvVar)
			} else {
				path = configService.GetFolderPath() + "/../locale"
			}

			log.Println("Loading locale files from " + path)

			return localize.NewSimpleLocalizer(ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger), apiSource, path), nil
		},
	}
}
