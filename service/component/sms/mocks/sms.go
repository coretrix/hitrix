package mocks

import (
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendOTPSMS(_ *beeorm.Engine, _ errorlogger.ErrorLogger, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendOTPCallout(_ *beeorm.Engine, _ errorlogger.ErrorLogger, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendMessage(_ *beeorm.Engine, _ errorlogger.ErrorLogger, message *sms.Message) error {
	return f.Called(message).Error(0)
}

func (f *FakeSMSSender) SendVerificationSMS(_ *beeorm.Engine, _ errorlogger.ErrorLogger, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendVerificationCallout(_ *beeorm.Engine, _ errorlogger.ErrorLogger, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) VerifyCode(_ *beeorm.Engine, _ errorlogger.ErrorLogger, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}
