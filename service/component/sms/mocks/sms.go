package mocks

import (
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/stretchr/testify/mock"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendOTPSMS(otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendOTPCallout(otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendMessage(message *sms.Message) error {
	return f.Called(message).Error(0)
}
