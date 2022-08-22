package registry

import (
	"errors"
	"log"
	"os"

	"github.com/sarulabs/di"
	"golang.org/x/oauth2/google"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	googlecalendar "github.com/coretrix/hitrix/service/component/google_calendar"
)

func ServiceProviderGoogleCalendar() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GoogleCalendarService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			credentialFile, ok := configService.String("google_calendar.credential_file")
			if !ok {
				return nil, errors.New("missing google_calendar.credential_file")
			}

			scopes, ok := configService.Strings("google_calendar.scopes")
			if !ok {
				return nil, errors.New("missing google_calendar.scopes")
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

			return googlecalendar.NewGoogleCalendar(oAuth2config), nil
		},
	}
}
