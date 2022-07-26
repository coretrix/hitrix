package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	googleanalytics "github.com/coretrix/hitrix/service/component/google_analytics"
)

func ServiceProviderGoogleAnalytics(newFunctions googleanalytics.NewProviderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleAnalyticsService,
		Build: func(ctn di.Container) (interface{}, error) {
			return googleanalytics.NewAPIManager(
				ctn.Get("config_directory").(string),
				ctn.Get(service.ConfigService).(config.IConfig),
				ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
				newFunctions)
		},
	}
}
