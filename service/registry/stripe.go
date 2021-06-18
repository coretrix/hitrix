package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	stripe2 "github.com/coretrix/hitrix/service/component/stripe"
	"github.com/sarulabs/di"
)

func ServiceDefinitionStripe() *service.Definition {
	return &service.Definition{
		Name:   service.StripeService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			conf := ctn.Get(service.ConfigService).(*config.Config).GetStringMap("stripe")
			appService := ctn.Get(service.AppService).(*app.App)

			var webhookSecrets map[string]string

			webhookVal, ok := conf["webhook_secrets"]
			if ok {
				webhookSecrets = webhookVal.(map[string]string)
			}

			return stripe2.NewStripe(conf["key"].(string), webhookSecrets, appService.Mode, conf), nil
		},
	}
}
