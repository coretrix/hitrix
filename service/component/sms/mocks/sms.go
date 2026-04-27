package mocks

import (
	"github.com/coretrix/trixorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/sms"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendMessage(_ *trixorm.Engine, message *sms.Message) error {
	return f.Called(message).Error(0)
}
