package mail

import (
	"encoding/json"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"

	"github.com/xorcare/pointer"

	"github.com/juju/errors"

	"github.com/latolukasz/beeorm"
	"github.com/mattbaird/gochimp"
)

type Mandrill struct {
	client           *gochimp.MandrillAPI
	defaultFromEmail string
	fromName         string
}

func NewMandrill(apiKey, defaultFromEmail, fromName string) *Mandrill {
	mandrillAPI, err := gochimp.NewMandrill(apiKey)

	if err != nil {
		panic(err)
	}

	return &Mandrill{client: mandrillAPI, defaultFromEmail: defaultFromEmail, fromName: fromName}
}

func (s *Mandrill) SendTemplate(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(ormService, message.From, message.To, message.Subject, message.TemplateName, message.TemplateData, nil, false)
}

func (s *Mandrill) SendTemplateAsync(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(ormService, message.From, message.To, message.Subject, message.TemplateName, message.TemplateData, nil, true)
}

func (s *Mandrill) SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error {
	return s.sendTemplate(ormService, message.From, message.To, message.Subject, message.TemplateName, message.TemplateData, message.Attachments, false)
}

func (s *Mandrill) SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *MessageAttachment) error {
	return s.sendTemplate(ormService, message.From, message.To, message.Subject, message.TemplateName, message.TemplateData, message.Attachments, true)
}

func (s *Mandrill) sendTemplate(ormService *beeorm.Engine, from string, to string, subject string, templateName string, templateData interface{}, attachments []gochimp.Attachment, async bool) error {
	if from == "" {
		from = s.defaultFromEmail
	}

	message := gochimp.Message{
		Subject:     subject,
		FromEmail:   from,
		Attachments: attachments,
		To: []gochimp.Recipient{
			{Email: to},
		},
	}

	mailTrackerEntity := &entity.MailTrackerEntity{
		Status:       entity.MailTrackerStatusNew,
		From:         from,
		To:           to,
		Subject:      subject,
		TemplateFile: templateName,
	}

	templateDataAsByte, err := json.Marshal(templateData)
	if err != nil {
		mailTrackerEntity.SenderError = "Cannot marshal TemplateData"
		mailTrackerEntity.Status = entity.MailTrackerStatusError
		ormService.Flush(mailTrackerEntity)

		return err
	}

	mailTrackerEntity.TemplateData = string(templateDataAsByte)

	var templateContent []gochimp.Var

	if templateData != nil {
		for key, value := range templateData.(map[string]interface{}) {
			templateContent = append(templateContent, *gochimp.NewVar(key, value))
		}
	}

	message.AddMergeVar(gochimp.MergeVars{Recipient: to, Vars: templateContent})
	responses, err := s.client.MessageSendTemplate(templateName, templateContent, message, async)

	if err != nil {
		mailTrackerEntity.SenderError = err.Error()
		mailTrackerEntity.Status = entity.MailTrackerStatusError
		ormService.Flush(mailTrackerEntity)

		return err
	}

	if responses != nil {
		for _, response := range responses {
			if response.RejectedReason != "" {
				mailTrackerEntity.SenderError += response.RejectedReason + " | "
			}
		}

		if mailTrackerEntity.SenderError != "" {
			mailTrackerEntity.Status = entity.MailTrackerStatusError
			ormService.Flush(mailTrackerEntity)
			return errors.New(mailTrackerEntity.SenderError)
		}
	}

	if async {
		mailTrackerEntity.Status = entity.MailTrackerStatusQueued
	} else {
		mailTrackerEntity.Status = entity.MailTrackerStatusSuccess
	}

	mailTrackerEntity.SentAt = pointer.Time(time.Now())
	ormService.Flush(mailTrackerEntity)

	return nil
}
