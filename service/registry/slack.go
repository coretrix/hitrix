package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/slack"
	"github.com/sarulabs/di"
)

func ServiceProviderSlack() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SlackService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			token, ok := configService.String("slack.token")
			if !ok {
				return nil, errors.New("missing slack.token")
			}

			errorChannel, _ := configService.String("slack.error_channel")
			devPanelURL, _ := configService.String("slack.dev_panel_url")

			return slack.NewSlackGo(token, errorChannel, devPanelURL), nil
		},
	}
}
