package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service/component/sms"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendMessage(_ *datalayer.DataLayer, message *sms.Message) error {
	return f.Called(message).Error(0)
}
