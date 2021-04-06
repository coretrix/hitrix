package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// FakeSysClock is a fake clock implementation
type FakeSysClock struct {
	mock.Mock
}

// Now is a no-op
func (c *FakeSysClock) Now() time.Time {
	return c.Called().Get(0).(time.Time)
}

// Now is a no-op
func (c *FakeSysClock) NowPointer() *time.Time {
	return c.Called().Get(0).(*time.Time)
}
