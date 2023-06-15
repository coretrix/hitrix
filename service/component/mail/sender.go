package mail

import (
	"encoding/json"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ISender interface {
	GetTemplateKeyFromConfig(templateName string) (string, error)
	SendTemplate(ormService *beeorm.Engine, message *Message) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error
	GetTemplateHTMLCode(templateName string) (string, error)
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

type Sender struct {
	ConfigService      config.IConfig
	ClockService       clock.IClock
	ErrorLoggerService errorlogger.ErrorLogger
	Provider           IProvider
}

func (s *Sender) GetTemplateKeyFromConfig(templateName string) (string, error) {
	return s.Provider.GetTemplateKeyFromConfig(s.ConfigService, templateName)
}

func (s *Sender) SendTemplate(ormService *beeorm.Engine, message *Message) error {
	if message.From == "" {
		message.From = s.Provider.GetDefaultFromEmail()
	}

	if message.FromName == "" {
		message.From = s.Provider.GetDefaultFromName()
	}

	mailTrackerEntity, err := s.createTrackingEntity(ormService, message)
	if err != nil {
		s.ErrorLoggerService.LogError(err)

		return err
	}

	err = s.Provider.SendTemplate(message)
	if err != nil {
		mailTrackerEntity.SenderError = err.Error()
		mailTrackerEntity.Status = entity.MailTrackerStatusError

		ormService.Flush(mailTrackerEntity)

		s.ErrorLoggerService.LogError(err)

		return err
	}

	mailTrackerEntity.Status = entity.MailTrackerStatusSuccess

	ormService.Flush(mailTrackerEntity)

	return nil
}

func (s *Sender) SendTemplateWithAttachments(ormService *beeorm.Engine, message *MessageAttachment) error {
	mailTrackerEntity, err := s.createTrackingEntity(ormService, &Message{
		From:         message.From,
		FromName:     message.FromName,
		ReplyTo:      message.ReplyTo,
		To:           message.To,
		Subject:      message.Subject,
		TemplateName: message.TemplateName,
		TemplateData: message.TemplateData,
	})
	if err != nil {
		s.ErrorLoggerService.LogError(err)

		return err
	}

	err = s.Provider.SendTemplateWithAttachments(message)
	if err != nil {
		mailTrackerEntity.SenderError = err.Error()
		mailTrackerEntity.Status = entity.MailTrackerStatusError

		ormService.Flush(mailTrackerEntity)

		s.ErrorLoggerService.LogError(err)

		return err
	}

	mailTrackerEntity.Status = entity.MailTrackerStatusSuccess

	ormService.Flush(mailTrackerEntity)

	return nil
}

func (s *Sender) GetTemplateHTMLCode(templateName string) (string, error) {
	return s.Provider.GetTemplateHTMLCode(templateName)
}

func (s *Sender) createTrackingEntity(ormService *beeorm.Engine, message *Message) (*entity.MailTrackerEntity, error) {
	mailTrackerEntity := &entity.MailTrackerEntity{
		Status:       entity.MailTrackerStatusNew,
		From:         message.From,
		To:           message.To,
		Subject:      message.Subject,
		TemplateFile: message.TemplateName,
		CreatedAt:    s.ClockService.Now(),
	}

	templateDataAsByte, err := json.Marshal(message.TemplateData)
	if err != nil {
		mailTrackerEntity.SenderError = err.Error()
		mailTrackerEntity.Status = entity.MailTrackerStatusError

		ormService.Flush(mailTrackerEntity)

		s.ErrorLoggerService.LogError(err)

		return mailTrackerEntity, err
	}

	mailTrackerEntity.TemplateData = string(templateDataAsByte)

	return mailTrackerEntity, nil
}
