package hitrix

import (
	"errors"

	"github.com/bluele/slack"
	"github.com/sarulabs/di"
)

func ServiceDefinitionSlackAPI() *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "slack_api",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			slackConfig := ctn.Get("config").(*Config).GetStringMap("slack")
			if slackConfig["token"] == "" {
				return nil, errors.New("missing slack.token")
			}
			if slackConfig["error_channel"] != "" && slackConfig["dev_panel_url"] != "" {
				return newSlack(slackConfig["token"].(string), slackConfig["error_channel"].(string), slackConfig["dev_panel_url"].(string)), nil
			}

			return newSlack(slackConfig["token"].(string), "", ""), nil
		},
	}
}

type SlackAPI struct {
	client       *slack.Slack
	errorChannel string
	devPanelURL  string
}

func newSlack(token, errorChannel, devPanelURL string) *SlackAPI {
	client := slack.New(token)

	return &SlackAPI{client: client, errorChannel: errorChannel, devPanelURL: devPanelURL}
}

func (s *SlackAPI) GetDevPanelURL() string {
	return s.devPanelURL
}

func (s *SlackAPI) GetErrorChannel() string {
	return s.errorChannel
}

func (s *SlackAPI) SendToChannel(channelName, message string, opt *slack.ChatPostMessageOpt) {
	err := s.client.ChatPostMessage(channelName, message, opt)
	if err != nil {
		panic(err)
	}
}
