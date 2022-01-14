package registry

import (
	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/slack"
)

func ServiceProviderErrorLogger() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ErrorLoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			slackAPIService, err := ctn.SafeGet(service.SlackService)

			if err == nil {
				return errorlogger.NewRedisErrorLogger(
					ctn.Get(service.AppService).(*app.App),
					ctn.Get(service.ORMEngineGlobalService).(*beeorm.Engine),
					slackAPIService.(slack.Slack),
				), nil
			}

			return errorlogger.NewRedisErrorLogger(
				ctn.Get(service.AppService).(*app.App),
				ctn.Get(service.ORMEngineGlobalService).(*beeorm.Engine),
				nil,
			), nil
		},
	}
}
