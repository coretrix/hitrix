package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/calendar"
)

func ServiceProviderMockCalendar(mock calendar.ICalendar) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.CalendarService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
