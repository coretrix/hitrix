package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/sentry"
)

func ServiceProviderSentry(tracesSampleRate *float64) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SentryService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			if configService == nil {
				panic("`config is nil")
			}

			backendVersionFinal := ""
			backendVersion, ok := configService.String("backend_version")
			if ok {
				backendVersionFinal = backendVersion
			}

			dsn, ok := configService.String("sentry.dsn")
			if !ok {
				panic("required config value sentry.dsn missing")
			}

			appService := ctn.Get(service.AppService).(*app.App)
			if appService == nil {
				panic("`app is nil")
			}

			return sentry.Init(dsn, appService.Mode, backendVersionFinal, tracesSampleRate), nil
		},
	}
}
