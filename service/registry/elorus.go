package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/elorus"
)

func ServiceProviderElorus() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ElorusService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			url, ok := configService.String("elorus.url")
			if !ok {
				return nil, errors.New("missing elorus.token key")
			}

			token, ok := configService.String("elorus.token")
			if !ok {
				return nil, errors.New("missing elorus.token key")
			}

			organizationId, ok := configService.String("elorus.organization_id")
			if !ok {
				return nil, errors.New("missing elorus.organization_id key")
			}

			appService := ctn.Get(service.AppService).(*app.App)

			return elorus.NewElorus(url, token, organizationId, appService.Mode), nil
		},
	}
}
