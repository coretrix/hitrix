package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

func ServiceProviderErrorLogger() *service.Definition {
	return &service.Definition{
		Name:   "error_logger",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			slackAPIService, err := ctn.SafeGet("slack_api")

			if err == nil {
				return errorlogger.NewRedisErrorLogger(
					ctn.Get("app").(*app.App),
					ctn.Get("orm_engine_global").(*orm.Engine),
					slackAPIService.(*slackapi.SlackAPI),
				), nil
			}

			return errorlogger.NewRedisErrorLogger(
				ctn.Get("app").(*app.App),
				ctn.Get("orm_engine_global").(*orm.Engine),
				nil,
			), nil
		},
	}
}
