package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/mattbaird/gochimp"
	"github.com/xorcare/pointer"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/config"
)

const (
	mandrillTemplateCachePrefix = "mandrill_template_"

	// ref: https://github.com/shawnmclean/Mandrill-dotnet/blob/05f26c917264751a903e3bcf83ca7153b5656526/src/Mandrill/Models/EmailMessage.cs#L19
	mergeLanguageHandlebars = "handlebars"
)

type Mandrill struct {
	client           *gochimp.MandrillAPI
	defaultFromEmail string
	fromName         string
}

func NewMandrill(configService config.IConfig) (Sender, error) {
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

	return &Mandrill{client: mandrillAPI, defaultFromEmail: fromEmail, fromName: fromName}, nil
}

func (s *Mandrill) GetTemplateKeyFromConfig(configService config.IConfig, templateName string) (string, error) {
	configPath := fmt.Sprintf("mail.mandrill.templates.%s", templateName)

	templateKey, ok := configService.String(configPath)
	if !ok {
		return "", fmt.Errorf("could not find email template key in config: %s", configPath)
	}

	return templateKey, nil
}

func (s *Mandrill) SendTemplate(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(
		ormService,
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		nil,
		false)
}

func (s *Mandrill) SendTemplateAsync(ormService *beeorm.Engine, message *Message) error {
	return s.sendTemplate(
		ormService,
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		nil,
		true)
}

func (s *Mandrill) SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error {
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
		ormService,
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		attachments,
		false)
}

func (s *Mandrill) SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *MessageAttachment) error {
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
		ormService,
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
		attachments,
		true)
}

func (s *Mandrill) sendTemplate(
	ormService *beeorm.Engine,
	from string,
	fromName string,
	to string,
	replyTo string,
	subject string,
	templateName string,
	templateData interface{},
	attachments []gochimp.Attachment,
	async bool,
) error {
	if from == "" {
		from = s.defaultFromEmail
	}

	message := gochimp.Message{
		MergeLanguage: mergeLanguageHandlebars,
		Subject:       subject,
		FromEmail:     from,
		Attachments:   attachments,
		To: []gochimp.Recipient{
			{Email: to},
		},
	}

	if fromName != "" {
		message.FromName = fromName
	}

	if replyTo != "" {
		message.Headers = map[string]string{
			"Reply-To": replyTo,
		}
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

func (s *Mandrill) GetTemplateHTMLCode(ormService *beeorm.Engine, templateName string, ttl int) (string, error) {
	key := mandrillTemplateCachePrefix + templateName
	redisCache := ormService.GetRedis()

	htmlCode, has := redisCache.Get(key)
	if has {
		return htmlCode, nil
	}

	template, err := s.client.TemplateInfo(templateName)
	if err != nil {
		return "", err
	}

	redisCache.Set(key, template.Code, ttl)

	return template.Code, nil
}
