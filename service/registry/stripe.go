package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/stripe"
)

func ServiceProviderStripe() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.StripeService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			key, _ := configService.String("stripe.key")

			secrets, _ := configService.StringMap("stripe.webhook_secrets")

			appService := ctn.Get(service.AppService).(*app.App)

			return stripe.NewStripe(key, secrets, appService), nil
		},
	}
}
