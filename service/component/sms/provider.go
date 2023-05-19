package sms

import (
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

const (
	success = "sent successfully"
	failure = "sent unsuccessfully"

	timeoutInSeconds = 5
)

type NewProviderFunc func(configService config.IConfig, clockService clock.IClock) (IProvider, error)

type IProvider interface {
	SendSMSMessage(msg *Message) (string, error)
	GetName() string
}

type Message struct {
	Text     string
	Number   string
	Provider *Provider
}

type Provider struct {
	Primary   IProvider
	Secondary IProvider
}
