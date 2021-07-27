package mocks

import (
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendOTPSMS(_ *beeorm.Engine, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendOTPCallout(_ *beeorm.Engine, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendMessage(_ *beeorm.Engine, message *sms.Message) error {
	return f.Called(message).Error(0)
}
