package mail

import (
	"github.com/latolukasz/beeorm"
	"github.com/mattbaird/gochimp"
)

type Sender interface {
	SendTemplate(ormService *beeorm.Engine, message *Message) error
	SendTemplateAsync(ormService *beeorm.Engine, message *Message) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error
	SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *MessageAttachment) error
	GetTemplateHTMLCode(ormService *beeorm.Engine, templateName string, ttl int) (string, error)
}

type Message struct {
	From         string
	To           string
	Subject      string
	TemplateName string
	TemplateData interface{}
}

type MessageAttachment struct {
	Message
	Attachments []gochimp.Attachment
}
