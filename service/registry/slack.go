package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/slack"
)

func ServiceProviderSlack() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SlackService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			botTokensConfig, ok := configService.Get("slack.bot_tokens")
			if !ok {
				return nil, errors.New("missing slack.bot_tokens")
			}

			botTokens := make(map[string]string)

			for name, token := range botTokensConfig.(map[interface{}]interface{}) {
				botTokens[name.(string)] = token.(string)
			}

			errorChannel, _ := configService.String("slack.error_channel")
			devPanelURL, _ := configService.String("slack.dev_panel_url")

			return slack.NewSlackGo(botTokens, errorChannel, devPanelURL), nil
		},
	}
}
