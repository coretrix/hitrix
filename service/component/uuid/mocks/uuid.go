package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockUUID struct {
	mock.Mock
}

func (m *MockUUID) Generate() string {
	return m.Called().String(0)
}
