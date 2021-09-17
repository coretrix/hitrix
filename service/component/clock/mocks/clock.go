package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type FakeSysClock struct {
	mock.Mock
}

func (c *FakeSysClock) Now() time.Time {
	return c.Called().Get(0).(time.Time)
}

func (c *FakeSysClock) NowPointer() *time.Time {
	return c.Called().Get(0).(*time.Time)
}
