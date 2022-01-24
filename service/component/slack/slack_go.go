package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

type APIClient struct {
	clients      map[string]*slack.Client
	errorChannel string
	devPanelURL  string
}

func NewSlackGo(botTokens map[string]string, errorChannel, devPanelURL string) *APIClient {
	clients := make(map[string]*slack.Client)
	for name, token := range botTokens {
		clients[name] = slack.New(token)
	}

	return &APIClient{clients: clients, errorChannel: errorChannel, devPanelURL: devPanelURL}
}

func (s *APIClient) GetDevPanelURL() string {
	return s.devPanelURL
}

func (s *APIClient) GetErrorChannel() string {
	return s.errorChannel
}

func (s *APIClient) SendToChannel(botName, channelName, message string, opt ...slack.MsgOption) error {
	client, ok := s.clients[botName]
	if !ok {
		return fmt.Errorf(`bot "%s" not defined`, botName)
	}

	opt = append(opt, slack.MsgOptionText(message, true))
	_, _, err := client.PostMessage(channelName, opt...)

	return err
}
