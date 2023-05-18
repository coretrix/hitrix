package mocks

import (
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/sms"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendMessage(_ *beeorm.Engine, _ errorlogger.ErrorLogger, message *sms.Message) error {
	return f.Called(message).Error(0)
}
