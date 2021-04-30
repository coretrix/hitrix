package slackapi

import (
	"github.com/slack-go/slack"
)

type SlackAPI struct {
	client       *slack.Client
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

func (s *SlackAPI) SendToChannel(channelName, message string, opt ...slack.MsgOption) {
	opt = append(opt, slack.MsgOptionText(message, true))
	_, _, err := s.client.PostMessage(channelName, opt...)
	if err != nil {
		panic(err)
	}
}
