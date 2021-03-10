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
			slackConfig := ctn.Get(service.ConfigService).(*config.Config).GetStringMap("slack")
			if slackConfig["token"] == "" {
				return nil, errors.New("missing slack.token")
			}
			if slackConfig["error_channel"] != "" && slackConfig["dev_panel_url"] != "" {
				return slackapi.NewSlack(slackConfig["token"].(string), slackConfig["error_channel"].(string), slackConfig["dev_panel_url"].(string)), nil
			}

			return slackapi.NewSlack(slackConfig["token"].(string), "", ""), nil
		},
	}
}
