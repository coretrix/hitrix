package mail

import (
	"errors"
	"fmt"

	"github.com/mattbaird/gochimp"

	"github.com/coretrix/hitrix/service/component/config"
)

const (
	// ref: https://github.com/shawnmclean/Mandrill-dotnet/blob/05f26c917264751a903e3bcf83ca7153b5656526/src/Mandrill/Models/EmailMessage.cs#L19
	mergeLanguageHandlebars = "handlebars"
)

type Mandrill struct {
	client           *gochimp.MandrillAPI
	defaultFromEmail string
	defaultFromName  string
}

func NewMandrill(configService config.IConfig) (IProvider, error) {
	apiKey, ok := configService.String("mail.mandrill.api_key")
	if !ok {
		return nil, errors.New("mail.mandrill.api_key is missing")
	}

	fromEmail, ok := configService.String("mail.mandrill.default_from_email")
	if !ok {
		return nil, errors.New("mail.mandrill.default_from_email is missing")
	}

	fromName, ok := configService.String("mail.mandrill.from_name")
	if !ok {
		return nil, errors.New("mail.mandrill.from_name is missing")
	}

	mandrillAPI, err := gochimp.NewMandrill(apiKey)

	if err != nil {
		panic(err)
	}

	return &Mandrill{client: mandrillAPI, defaultFromEmail: fromEmail, defaultFromName: fromName}, nil
}

func (s *Mandrill) GetTemplateKeyFromConfig(configService config.IConfig, templateName string) (string, error) {
	configPath := fmt.Sprintf("mail.mandrill.templates.%s", templateName)

	templateKey, ok := configService.String(configPath)
	if !ok {
		return "", fmt.Errorf("could not find email template key in config: %s", configPath)
	}

	return templateKey, nil
}

func (s *Mandrill) SendTemplate(message *Message) error {
	return s.sendTemplate(
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		nil)
}

func (s *Mandrill) SendTemplateWithAttachments(message *MessageAttachment) error {
	var attachments []gochimp.Attachment
	if message.Attachments != nil {
		attachments = make([]gochimp.Attachment, len(message.Attachments))
		for i, attachment := range message.Attachments {
			attachments[i] = gochimp.Attachment{
				Type:    attachment.ContentType,
				Name:    attachment.Filename,
				Content: attachment.Base64Content,
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
		attachments)
}

func (s *Mandrill) sendTemplate(
	from string,
	fromName string,
	to string,
	replyTo string,
	subject string,
	templateName string,
	templateData interface{},
	attachments []gochimp.Attachment,
) error {
	message := gochimp.Message{
		MergeLanguage: mergeLanguageHandlebars,
		Subject:       subject,
		FromName:      fromName,
		FromEmail:     from,
		Attachments:   attachments,
		To: []gochimp.Recipient{
			{Email: to},
		},
	}

	if replyTo != "" {
		message.Headers = map[string]string{
			"Reply-To": replyTo,
		}
	}

	var templateContent []gochimp.Var

	if templateData != nil {
		for key, value := range templateData.(map[string]interface{}) {
			templateContent = append(templateContent, *gochimp.NewVar(key, value))
		}
	}

	message.AddMergeVar(gochimp.MergeVars{Recipient: to, Vars: templateContent})
	responses, err := s.client.MessageSendTemplate(templateName, templateContent, message, false)

	if err != nil {
		return err
	}

	if responses != nil {
		var errAsString string

		for _, response := range responses {
			if response.RejectedReason != "" {
				errAsString += response.RejectedReason + " | "
			}
		}

		if errAsString != "" {
			return errors.New(errAsString)
		}
	}

	return nil
}

func (s *Mandrill) GetTemplateHTMLCode(templateName string) (string, error) {
	template, err := s.client.TemplateInfo(templateName)
	if err != nil {
		return "", err
	}

	return template.Code, nil
}

func (s *Mandrill) GetDefaultFromEmail() string {
	return s.defaultFromEmail
}

func (s *Mandrill) GetDefaultFromName() string {
	return s.defaultFromName
}
