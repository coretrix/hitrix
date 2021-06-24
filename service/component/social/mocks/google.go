package mocks

import (
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/stretchr/testify/mock"
)

type Google struct {
	mock.Mock
}

func (m *Google) GetUserData(token string) social.IUserData {
	return m.Called(token).Get(0).(social.IUserData)
}
