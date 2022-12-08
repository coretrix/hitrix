package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
)

func ServiceProviderRequestLogger() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.RequestLoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := service.DI().OrmConfig()

			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.RequestLoggerEntity"]; !ok {
				panic("you should register RequestLoggerEntity")
			}

			return requestlogger.NewDBLogger(ctn.Get(service.ClockService).(clock.IClock)), nil
		},
	}
}
