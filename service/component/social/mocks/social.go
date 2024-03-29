package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/social"
)

type Social struct {
	mock.Mock
}

func (m *Social) GetUserData(_ context.Context, token string, isAndroid bool) (*social.UserData, error) {
	args := m.Called(token, isAndroid)

	return args.Get(0).(*social.UserData), args.Error(1)
}
