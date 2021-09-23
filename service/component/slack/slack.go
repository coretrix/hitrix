package slack

import (
	"github.com/slack-go/slack"
)

type Slack interface {
	GetDevPanelURL() string
	GetErrorChannel() string
	SendToChannel(channelName, message string, opt ...slack.MsgOption) error
}

//TODO refactor
