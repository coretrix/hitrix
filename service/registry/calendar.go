package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/calendar"
	"github.com/coretrix/hitrix/service/component/config"
)

func ServiceProviderCalendar(newFunc calendar.NewCalendarFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CalendarService,
		Build: func(ctn di.Container) (interface{}, error) {
			return newFunc(ctn.Get(service.ConfigService).(config.IConfig))
		},
	}
}
