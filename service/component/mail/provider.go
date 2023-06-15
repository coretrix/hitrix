package mail

import (
	"github.com/coretrix/hitrix/service/component/config"
)

type NewSenderFunc func(configService config.IConfig) (IProvider, error)

type IProvider interface {
	GetTemplateKeyFromConfig(configService config.IConfig, templateName string) (string, error)
	SendTemplate(message *Message) error
	GetDefaultFromEmail() string
	GetDefaultFromName() string
	SendTemplateWithAttachments(message *MessageAttachment) error
	GetTemplateHTMLCode(templateName string) (string, error)
}
