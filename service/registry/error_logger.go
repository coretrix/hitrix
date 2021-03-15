package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

func ServiceProviderErrorLogger() *service.Definition {
	return &service.Definition{
		Name:   service.ErrorLoggerService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			slackAPIService, err := ctn.SafeGet(service.SlackAPIService)

			if err == nil {
				return errorlogger.NewRedisErrorLogger(
					ctn.Get(service.AppService).(*app.App),
					ctn.Get(service.ORMEngineGlobalService).(*orm.Engine),
					slackAPIService.(*slackapi.SlackAPI),
				), nil
			}

			return errorlogger.NewRedisErrorLogger(
				ctn.Get(service.AppService).(*app.App),
				ctn.Get(service.ORMEngineGlobalService).(*orm.Engine),
				nil,
			), nil
		},
	}
}
