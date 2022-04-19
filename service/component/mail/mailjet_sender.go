package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/mailjet/mailjet-apiv3-go/resources"
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/xorcare/pointer"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/config"
)

const (
	mailjetTemplateCachePrefix = "mailjet_template_"
)

type Mailjet struct {
	client           *mailjet.Client
	defaultFromEmail string
	fromName         string
	sandboxMode      bool
}

func NewMailjet(configService config.IConfig) (Sender, error) {
	apiKeyPublic, ok := configService.String("mail.mailjet.api_key_public")
	if !ok {
		return nil, errors.New("mail.mailjet.api_key_public is missing")
	}

	apiKeyPrivate, ok := configService.String("mail.mailjet.api_key_private")
	if !ok {
		return nil, errors.New("mail.mailjet.api_key_private is missing")
	}

	fromEmail, ok := configService.String("mail.mailjet.default_from_email")
	if !ok {
		return nil, errors.New("mail.mailjet.default_from_email is missing")
	}

	fromName, ok := configService.String("mail.mailjet.from_name")
	if !ok {
		return nil, errors.New("mail.mailjet.from_name is missing")
	}

	sandboxMode, ok := configService.Bool("mail.mailjet.sandbox_mode")
	if !ok {
		return nil, errors.New("mail.mailjet.sandbox_mode is missing")
	}

	mailjetAPI := mailjet.NewMailjetClient(apiKeyPublic, apiKeyPrivate)

	return &Mailjet{client: mailjetAPI, defaultFromEmail: fromEmail, fromName: fromName, sandboxMode: sandboxMode}, nil
}

func (s *Mailjet) SendTemplate(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(ormService, message.From, message.FromName, message.To, message.ReplyTo, message.Subject, message.TemplateName, message.TemplateData, nil, false)
}

func (s *Mailjet) SendTemplateAsync(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(ormService, message.From, message.FromName, message.To, message.ReplyTo, message.Subject, message.TemplateName, message.TemplateData, nil, true)
}

func (s *Mailjet) SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error {
	var attachments []mailjet.AttachmentV31
	if message.Attachments != nil {
		attachments = make([]mailjet.AttachmentV31, len(message.Attachments))
		for i, attachment := range message.Attachments {
			attachments[i] = mailjet.AttachmentV31{
				ContentType:   attachment.ContentType,
				Filename:      attachment.Filename,
				Base64Content: attachment.Base64Content,
			}
		}
	}

	return s.sendTemplate(ormService, message.From, message.FromName, message.To, message.ReplyTo, message.Subject, message.TemplateName, message.TemplateData, attachments, false)
}

func (s *Mailjet) SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *MessageAttachment) error {
	var attachments []mailjet.AttachmentV31
	if message.Attachments != nil {
		attachments = make([]mailjet.AttachmentV31, len(message.Attachments))
		for i, attachment := range message.Attachments {
			attachments[i] = mailjet.AttachmentV31{
				ContentType:   attachment.ContentType,
				Filename:      attachment.Filename,
				Base64Content: attachment.Base64Content,
			}
		}
	}

	return s.sendTemplate(ormService, message.From, message.FromName, message.To, message.ReplyTo, message.Subject, message.TemplateName, message.TemplateData, attachments, true)
}

func (s *Mailjet) sendTemplate(ormService *beeorm.Engine, from string, fromName string, to string, replyTo string, subject string, templateName string, templateData interface{}, attachments []mailjet.AttachmentV31, async bool) error {
	if from == "" {
		from = s.defaultFromEmail
	}

	templateID, err := strconv.ParseInt(templateName, 10, 64)
	if err != nil {
		return err
	}

	messageInfo := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: from,
			Name:  fromName,
		},
		To: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: to,
			},
		},
		Subject:          subject,
		Variables:        templateData.(map[string]interface{}),
		TemplateID:       templateID,
		TemplateLanguage: true,
	}

	if len(replyTo) > 0 {
		messageInfo.ReplyTo = &mailjet.RecipientV31{
			Email: replyTo,
		}
	}

	if len(attachments) > 0 {
		messageInfo.Attachments = (*mailjet.AttachmentsV31)(&attachments)
	}

	message := &mailjet.MessagesV31{
		Info:        []mailjet.InfoMessagesV31{messageInfo},
		SandBoxMode: s.sandboxMode,
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

	results, err := s.client.SendMailV31(message)

	if err != nil {
		mailTrackerEntity.SenderError = err.Error()
		mailTrackerEntity.Status = entity.MailTrackerStatusError
		ormService.Flush(mailTrackerEntity)

		return err
	}

	if results != nil {
		for _, response := range results.ResultsV31 {
			if response.Status != "success" {
				mailTrackerEntity.SenderError += "error | "
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

func (s *Mailjet) GetTemplateHTMLCode(ormService *beeorm.Engine, templateName string, ttl int) (string, error) {
	key := mailjetTemplateCachePrefix + templateName
	redisCache := ormService.GetRedis()

	var templates []resources.TemplateDetailcontent
	err := s.client.Get(&mailjet.Request{
		Resource: "template",
		AltID:    templateName,
		Action:   "detailcontent",
	}, &templates)
	if err != nil {
		return "", err
	}

	if len(templates) == 0 {
		return "", fmt.Errorf("no template found with name %v", templateName)
	}

	if len(templates) > 1 {
		return "", fmt.Errorf("%d templates found with name %v", len(templates), templateName)
	}

	template := templates[0]

	redisCache.Set(key, template.HtmlPart, ttl)

	return template.HtmlPart, nil
}
