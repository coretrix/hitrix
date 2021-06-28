package mocks

import (
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/latolukasz/orm"
	"github.com/stretchr/testify/mock"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendOTPSMS(ormService *orm.Engine, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendOTPCallout(ormService *orm.Engine, otp *sms.OTP) error {
	return f.Called(otp).Error(0)
}

func (f *FakeSMSSender) SendMessage(ormService *orm.Engine, message *sms.Message) error {
	return f.Called(message).Error(0)
}
