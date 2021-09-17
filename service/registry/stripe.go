package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	stripe2 "github.com/coretrix/hitrix/service/component/stripe"
	"github.com/sarulabs/di"
)

func ServiceProviderStripe() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.StripeService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			key, ok := configService.String("stripe.key")
			if ok {
				return nil, errors.New("missing stripe key")
			}

			secrets, ok := configService.StringMap("stripe.webhook_secrets")
			if !ok {
				return nil, errors.New("missing stripe secrets")
			}

			appService := ctn.Get(service.AppService).(*app.App)

			return stripe2.NewStripe(key, secrets, appService.Mode), nil
		},
	}
}
