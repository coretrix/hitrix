package mocks

import (
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/latolukasz/orm"
	"github.com/stretchr/testify/mock"
)

type Sender struct {
	mock.Mock
}

func (m *Sender) SendTemplate(_ *orm.Engine, message *mail.TemplateMessage) error {
	return m.Called(message).Error(0)
}

func (m *Sender) SendTemplateAsync(_ *orm.Engine, message *mail.TemplateMessage) error {
	return m.Called(message).Error(0)
}

func (m *Sender) SendTemplateWithAttachments(_ *orm.Engine, message *mail.TemplateAttachmentMessage) error {
	return m.Called(message).Error(0)
}

func (m *Sender) SendTemplateWithAttachmentsAsync(_ *orm.Engine, message *mail.TemplateAttachmentMessage) error {
	return m.Called(message).Error(0)
}
