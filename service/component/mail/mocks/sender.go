package mocks

import (
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/mail"
)

type Sender struct {
	mock.Mock
}

func (m *Sender) SendTemplate(_ *beeorm.Engine, message *mail.Message) error {
	return m.Called(message.To).Error(0)
}

func (m *Sender) SendTemplateAsync(_ *beeorm.Engine, message *mail.Message) error {
	return m.Called(message.To).Error(0)
}

func (m *Sender) SendTemplateWithAttachments(_ *beeorm.Engine, message *mail.MessageAttachment) error {
	return m.Called(message.To).Error(0)
}

func (m *Sender) SendTemplateWithAttachmentsAsync(_ *beeorm.Engine, message *mail.MessageAttachment) error {
	return m.Called(message.To).Error(0)
}

func (m *Sender) GetTemplateHTMLCode(_ *beeorm.Engine, templateName string, ttl int) (string, error) {
	args := m.Called(templateName, ttl)
	return args.Get(0).(string), args.Error(1)
}
