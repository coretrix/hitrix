package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	googlecalendar "github.com/coretrix/hitrix/service/component/google_calendar"
)

func ServiceProviderMockGoogleCalendar(mock googlecalendar.IGoogleCalendar) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleCalendarService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
