package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"
	"github.com/juju/errors"
	"github.com/sarulabs/di"
)

func ServiceDefinitionSlackAPI() *service.Definition {
	return &service.Definition{
		Name:   service.SlackAPIService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			token, ok := configService.String("slack.token")
			if !ok {
				return nil, errors.New("missing slack.token")
			}

			errorChannel, _ := configService.String("slack.error_channel")
			devPanelURL, _ := configService.String("slack.dev_panel_url")

			return slackapi.NewSlack(token, errorChannel, devPanelURL), nil
		},
	}
}
