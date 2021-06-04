package mail

import (
	"github.com/latolukasz/orm"
	"github.com/mattbaird/gochimp"
)

type Sender interface {
	SendTemplate(ormService *orm.Engine, message *Message) error
	SendTemplateAsync(ormService *orm.Engine, message *Message) error
	SendTemplateWithAttachments(ormService *orm.Engine, message *MessageAttachment) error
	SendTemplateWithAttachmentsAsync(ormService *orm.Engine, message *MessageAttachment) error
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
