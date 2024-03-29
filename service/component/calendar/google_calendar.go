package calendar

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/service/component/config"
)

func NewGoogleCalendar(configService config.IConfig) (ICalendar, error) {
	credentialFile, ok := configService.String("calendar.credential_file")
	if !ok {
		return nil, errors.New("missing calendar.credential_file")
	}

	scopes, ok := configService.Strings("calendar.scopes")
	if !ok {
		return nil, errors.New("missing calendar.scopes")
	}

	if len(scopes) == 0 {
		return nil, errors.New("google calendar, scopes list is empty")
	}

	credentials, err := os.ReadFile(credentialFile)
	if err != nil {
		return nil, errors.New("google calendar, unable to read client secret file")
	}

	oAuth2config, err := google.ConfigFromJSON(credentials, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return &GoogleCalendar{
		Oauth2Config: oAuth2config,
		Ctx:          context.Background(),
	}, nil
}

type GoogleCalendar struct {
	Oauth2Config *oauth2.Config
	Ctx          context.Context
}

func (gc *GoogleCalendar) GetAuthLink(state string) string {
	// we're forcing approval since if we don't google won't return refresh token second time this user logs in
	return gc.Oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (gc *GoogleCalendar) GetStateCodeFromGin(c *gin.Context) (string, string, error) {
	state, hasState := c.GetQuery("state")
	code, hasCode := c.GetQuery("code")

	if hasState && hasCode && state != "" && code != "" {
		return state, code, nil
	}

	return "", "", errors.New("google calendar API, failed to obtain code or/and state from context")
}

func (gc *GoogleCalendar) GetTokenFromCode(code string) (*oauth2.Token, error) {
	tok, err := gc.Oauth2Config.Exchange(gc.Ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google calendar API, unable to retrieve token from web: %v", err)
	}

	return tok, nil
}

// RefreshToken refreshes the token if it's expired otherwise it returns the same token
// useful for updating token inside your database when it's refreshed
func (gc *GoogleCalendar) RefreshToken(token *oauth2.Token) (bool, *oauth2.Token, error) {
	tokenSource := gc.Oauth2Config.TokenSource(gc.Ctx, token)

	newToken, err := tokenSource.Token()

	if err != nil {
		return false, nil, err
	}

	return newToken.AccessToken != token.AccessToken, newToken, nil
}

func (gc *GoogleCalendar) GetCalendars(token *oauth2.Token) ([]*calendar.CalendarListEntry, error) {
	srv, err := calendar.NewService(gc.Ctx, option.WithHTTPClient(gc.Oauth2Config.Client(gc.Ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("google calendar API, unable to retrieve Calendar client: %v", err)
	}

	calendars := make([]*calendar.CalendarListEntry, 0)

	err = srv.CalendarList.List().Pages(gc.Ctx, func(list *calendar.CalendarList) error {
		calendars = append(calendars, list.Items...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return calendars, nil
}

type GetCalendarEventsArgs struct {
	ShowDeleted  *bool
	SingleEvents *bool
	TimeZone     string
	TimeMin      string
	TimeMax      string
	UpdatedMin   string
}

func (gc *GoogleCalendar) GetCalendarEvents(token *oauth2.Token, calendarID string, args *GetCalendarEventsArgs) ([]*calendar.Event, error) {
	srv, err := calendar.NewService(gc.Ctx, option.WithHTTPClient(gc.Oauth2Config.Client(gc.Ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("google calendar API, unable to retrieve Calendar client: %v", err)
	}

	eventsConfig := srv.Events.List(calendarID)

	if args.ShowDeleted != nil {
		eventsConfig.ShowDeleted(*args.ShowDeleted)
	}

	if args.SingleEvents != nil {
		eventsConfig.SingleEvents(*args.SingleEvents)
	}

	if args.TimeMax != "" {
		eventsConfig.TimeMax(args.TimeMax)
	}

	if args.TimeMin != "" {
		eventsConfig.TimeMin(args.TimeMin)
	}

	if args.TimeZone != "" {
		eventsConfig.TimeZone(args.TimeZone)
	}

	if args.UpdatedMin != "" {
		eventsConfig.UpdatedMin(args.UpdatedMin)
	}

	eventsRes, err := eventsConfig.Do()

	if err != nil {
		return nil, fmt.Errorf("google calendar API, unable to retrieve Calendar events list: %v", err)
	}

	return eventsRes.Items, nil
}

func (gc *GoogleCalendar) UpsertEvent(token *oauth2.Token, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	srv, err := calendar.NewService(gc.Ctx, option.WithHTTPClient(gc.Oauth2Config.Client(gc.Ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("google calendar API, unable to retrieve Calendar client: %v", err)
	}

	if event.Created == "" {
		createdEvent, err := srv.Events.Insert(calendarID, event).Do()

		return createdEvent, err
	}

	updatedEvent, err := srv.Events.Update(calendarID, event.Id, event).Do()

	return updatedEvent, err
}
