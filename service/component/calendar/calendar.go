package calendar

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"

	"github.com/coretrix/hitrix/service/component/config"
)

type NewCalendarFunc func(configService config.IConfig) (ICalendar, error)

type ICalendar interface {
	GetAuthLink(state string) string
	GetStateCodeFromGin(c *gin.Context) (string, string, error)
	GetTokenFromCode(code string) (*oauth2.Token, error)
	RefreshToken(token *oauth2.Token) (bool, *oauth2.Token, error)
	GetCalendars(token *oauth2.Token) ([]*calendar.CalendarListEntry, error)
	GetCalendarEvents(token *oauth2.Token, calendarID string, args *GetCalendarEventsArgs) ([]*calendar.Event, error)
	UpsertEvent(token *oauth2.Token, calendarID string, event *calendar.Event) (*calendar.Event, error)
}
