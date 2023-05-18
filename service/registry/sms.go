package registry

import (
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/sms"
)

func ServiceProviderSMS(entity sms.LogEntity, newFuncProviders ...sms.NewProviderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SMSService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			clockService := ctn.Get(service.ClockService).(clock.IClock)

			providerContainer := map[string]sms.IProvider{}
			for _, newFunc := range newFuncProviders {
				provider, err := newFunc(configService, clockService)
				if err != nil {
					panic(err)
				}

				providerContainer[provider.GetName()] = provider
			}

			return &sms.Sender{
				Logger:            entity,
				Clock:             clockService,
				ProviderContainer: providerContainer,
			}, nil
		},
	}
}
