package mocks

import (
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/stretchr/testify/mock"
)

type Facebook struct {
	mock.Mock
}

func (m *Facebook) GetUserData(token string) social.IUserData {
	return m.Called(token).Get(0).(social.IUserData)
}
