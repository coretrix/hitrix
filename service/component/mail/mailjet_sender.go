package mail

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mailjet/mailjet-apiv3-go/resources"
	"github.com/mailjet/mailjet-apiv3-go/v3"

	"github.com/coretrix/hitrix/service/component/config"
)

type Mailjet struct {
	client           *mailjet.Client
	defaultFromEmail string
	defaultFromName  string
	sandboxMode      bool
}

func NewMailjet(configService config.IConfig) (IProvider, error) {
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

	fromName, ok := configService.String("mail.mailjet.default_from_name")
	if !ok {
		return nil, errors.New("mail.mailjet.from_name is missing")
	}

	sandboxMode, ok := configService.Bool("mail.mailjet.sandbox_mode")
	if !ok {
		return nil, errors.New("mail.mailjet.sandbox_mode is missing")
	}

	mailjetAPI := mailjet.NewMailjetClient(apiKeyPublic, apiKeyPrivate)

	return &Mailjet{client: mailjetAPI, sandboxMode: sandboxMode, defaultFromEmail: fromEmail, defaultFromName: fromName}, nil
}

func (s *Mailjet) GetTemplateKeyFromConfig(configService config.IConfig, templateName string) (string, error) {
	configPath := fmt.Sprintf("mail.mailjet.templates.%s", templateName)

	templateKey, ok := configService.String(configPath)
	if !ok {
		return "", fmt.Errorf("could not find email template key in config: %s", configPath)
	}

	return templateKey, nil
}

func (s *Mailjet) SendTemplate(message *Message) error {
	return s.sendTemplate(
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		nil,
	)
}

func (s *Mailjet) SendTemplateWithAttachments(message *MessageAttachment) error {
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

	return s.sendTemplate(
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		attachments,
	)
}

func (s *Mailjet) sendTemplate(
	from string,
	fromName string,
	to string,
	replyTo string,
	subject string,
	templateName string,
	templateData interface{},
	attachments []mailjet.AttachmentV31,
) error {
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

	results, err := s.client.SendMailV31(message)

	if err != nil {
		return err
	}

	if results != nil {
		var errAsString string

		for _, response := range results.ResultsV31 {
			if response.Status != "success" {
				errAsString += "error | "
			}
		}

		if errAsString != "" {
			return errors.New(errAsString)
		}
	}

	return nil
}

func (s *Mailjet) GetTemplateHTMLCode(templateName string) (string, error) {
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

	return template.HtmlPart, nil
}

func (s *Mailjet) GetDefaultFromEmail() string {
	return s.defaultFromEmail
}

func (s *Mailjet) GetDefaultFromName() string {
	return s.defaultFromName
}
