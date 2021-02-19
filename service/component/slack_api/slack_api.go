package slackapi

import (
	"github.com/bluele/slack"
)

type SlackAPI struct {
	client       *slack.Slack
	errorChannel string
	devPanelURL  string
}

func NewSlack(token, errorChannel, devPanelURL string) *SlackAPI {
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
