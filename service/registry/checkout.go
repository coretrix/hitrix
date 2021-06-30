package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/checkout"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"
	"github.com/xorcare/pointer"
)

func ServiceDefinitionCheckout() *service.Definition {
	return &service.Definition{
		Name:   service.CheckoutService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			secretKey, ok := configService.String("checkout.secret_key")
			if ok {
				return nil, errors.New("missing checkout.secret_key key")
			}

			webhookSecrets, ok := configService.StringMap("checkout.webhook_keys")
			if ok {
				return nil, errors.New("missing checkout.webhook_keys key")
			}

			var publicKey *string

			if publicKeyVal, ok := configService.String("checkout.public_key"); ok && publicKeyVal != "" {
				publicKey = pointer.String(publicKeyVal)
			}

			appService := ctn.Get(service.AppService).(*app.App)

			return checkout.NewCheckout(secretKey, publicKey, appService.Mode, webhookSecrets), nil
		},
	}
}
