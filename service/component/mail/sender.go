package mail

import (
	"github.com/latolukasz/orm"
	"github.com/mattbaird/gochimp"
)

type Sender interface {
	SendTemplate(ormService *orm.Engine, message *TemplateMessage) error
	SendTemplateAsync(ormService *orm.Engine, message *TemplateMessage) error
	SendTemplateWithAttachments(ormService *orm.Engine, message *TemplateAttachmentMessage) error
	SendTemplateWithAttachmentsAsync(ormService *orm.Engine, message *TemplateAttachmentMessage) error
}

type TemplateMessage struct {
	From         string
	To           string
	Subject      string
	templateName string
	templateData interface{}
}

type TemplateAttachmentMessage struct {
	TemplateMessage
	Attachments []gochimp.Attachment
}
