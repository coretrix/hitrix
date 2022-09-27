package mail

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service/component/config"
)

type NewSenderFunc func(configService config.IConfig) (Sender, error)

type Sender interface {
	GetTemplateKeyFromConfig(configService config.IConfig, templateName string) (string, error)
	SendTemplate(ormService *beeorm.Engine, message *Message) error
	SendTemplateAsync(ormService *beeorm.Engine, message *Message) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error
	SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *MessageAttachment) error
	GetTemplateHTMLCode(ormService *beeorm.Engine, templateName string, ttl int) (string, error)
}

type Message struct {
	From         string
	FromName     string
	ReplyTo      string
	To           string
	Subject      string
	TemplateName string
	TemplateData interface{}
}

type Attachment struct {
	ContentType   string
	Filename      string
	Base64Content string
}

type MessageAttachment struct {
	Message
	Attachments []Attachment
}
