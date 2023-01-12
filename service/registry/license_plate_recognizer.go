package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	licenseplaterecognizer "github.com/coretrix/hitrix/service/component/license_plate_recognizer"
)

func ServiceProviderLicensePlateRecognizer() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LicensePlateRecognizerService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			apiKey, ok := configService.String("platerecognizer.api_key")
			if !ok {
				return nil, errors.New("missing platerecognizer.api_key")
			}

			return licenseplaterecognizer.NewPlateRecognizer(apiKey), nil
		},
	}
}
