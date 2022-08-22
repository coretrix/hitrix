package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"

	hitrixcalendar "github.com/coretrix/hitrix/service/component/calendar"
)

type FakeCalendar struct {
	mock.Mock
}

func (gc *FakeCalendar) GetAuthLink(state string) string {
	args := gc.Called(state)

	return args.Get(0).(string)
}

func (gc *FakeCalendar) GetStateCodeFromGin(c *gin.Context) (string, string, error) {
	args := gc.Called(c)

	return args.Get(0).(string), args.Get(1).(string), args.Error(2)
}

func (gc *FakeCalendar) GetTokenFromCode(code string) (*oauth2.Token, error) {
	args := gc.Called(code)

	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (gc *FakeCalendar) RefreshToken(token *oauth2.Token) (bool, *oauth2.Token, error) {
	args := gc.Called(token)

	return args.Get(0).(bool), args.Get(1).(*oauth2.Token), args.Error(2)
}

func (gc *FakeCalendar) GetCalendars(token *oauth2.Token) ([]*calendar.CalendarListEntry, error) {
	args := gc.Called(token)

	return args.Get(0).([]*calendar.CalendarListEntry), args.Error(1)
}

func (gc *FakeCalendar) GetCalendarEvents(token *oauth2.Token,
	calendarID string,
	args *hitrixcalendar.GetCalendarEventsArgs) ([]*calendar.Event, error) {
	args2 := gc.Called(token, calendarID, args)

	return args2.Get(0).([]*calendar.Event), args2.Error(1)
}

func (gc *FakeCalendar) UpsertEvent(token *oauth2.Token, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	args := gc.Called(token, calendarID, event)

	return args.Get(0).(*calendar.Event), args.Error(1)
}
