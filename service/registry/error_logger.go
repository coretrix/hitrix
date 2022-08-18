package registry

import (
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/sentry"
	"github.com/coretrix/hitrix/service/component/slack"
)

func ServiceProviderErrorLogger() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ErrorLoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			var sentryService sentry.ISentry = nil
			var slackAPIService slack.Slack = nil

			sentryServiceInterface, err := ctn.SafeGet(service.SentryService)
			if err != nil {
				sentryService = sentryServiceInterface.(sentry.ISentry)
			}

			slackAPIServiceInterface, err := ctn.SafeGet(service.SlackService)
			if err != nil {
				slackAPIService = slackAPIServiceInterface.(slack.Slack)
			}

			return errorlogger.NewRedisErrorLogger(
				ctn.Get(service.AppService).(*app.App),
				ctn.Get(service.ORMEngineGlobalService).(*beeorm.Engine),
				slackAPIService,
				sentryService,
			), nil
		},
	}
}
