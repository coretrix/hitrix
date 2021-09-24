package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/trustpilot"
	"github.com/juju/errors"
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"
)

func ServiceDefinitionTrustpilot() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.TrustpilotService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			apiKey, ok := configService.String("trustpilot.apiKey")
			if !ok {
				return nil, errors.New("missing trustpilot.apiKey")
			}

			apiSecret, ok := configService.String("trustpilot.apiSecret")
			if !ok {
				return nil, errors.New("missing trustpilot.apiSecret")
			}

			username, ok := configService.String("trustpilot.username")
			if !ok {
				return nil, errors.New("missing trustpilot.username")
			}

			password, ok := configService.String("trustpilot.password")
			if !ok {
				return nil, errors.New("missing trustpilot.password")
			}

			businessUnitID, ok := configService.String("trustpilot.businessUnitID")
			if !ok {
				return nil, errors.New("missing trustpilot.businessUnitID")
			}

			ormService := ctn.Get(service.ORMEngineGlobalService).(*beeorm.Engine)

			clockService := ctn.Get(service.ClockService).(clock.IClock)

			return trustpilot.NewTrustpilot(apiKey, apiSecret, username, password, businessUnitID, ormService, clockService)
		},
	}
}
