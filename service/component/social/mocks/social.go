package mocks

import (
	"context"

	"github.com/coretrix/hitrix/service/component/social"
	"github.com/stretchr/testify/mock"
)

type Social struct {
	mock.Mock
}

func (m *Social) GetUserData(_ context.Context, token string) (*social.UserData, error) {
	return m.Called(token).Get(0).(*social.UserData), m.Called(token).Error(1)
}

func (m *Social) SetIsAndroid(_ bool) {}
