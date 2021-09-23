package slack

import "github.com/slack-go/slack"

type APIClient struct {
	client       *slack.Client
	errorChannel string
	devPanelURL  string
}

func NewSlackGo(token, errorChannel, devPanelURL string) *APIClient {
	client := slack.New(token)

	return &APIClient{client: client, errorChannel: errorChannel, devPanelURL: devPanelURL}
}

func (s *APIClient) GetDevPanelURL() string {
	return s.devPanelURL
}

func (s *APIClient) GetErrorChannel() string {
	return s.errorChannel
}

func (s *APIClient) SendToChannel(channelName, message string, opt ...slack.MsgOption) error {
	opt = append(opt, slack.MsgOptionText(message, true))
	_, _, err := s.client.PostMessage(channelName, opt...)

	return err
}
