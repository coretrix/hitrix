package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/goroutine"
)

func ServiceProviderGoroutine() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoroutineService,
		Build: func(ctn di.Container) (interface{}, error) {
			return goroutine.NewGoroutineManager(ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)), nil
		},
	}
}
