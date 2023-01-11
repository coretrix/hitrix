package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	licenceplaterecognizer "github.com/coretrix/hitrix/service/component/license_plate_recognizer"
)

func ServiceProviderLicencePlateRecognizer() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.LicencePlateRecognizerService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			apiKey, ok := configService.String("platerecognizer.api_key")
			if !ok {
				return nil, errors.New("missing platerecognizer.api_key")
			}

			return licenceplaterecognizer.NewPlateRecognizer(apiKey), nil
		},
	}
}
