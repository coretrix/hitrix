package mocks

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
)

type FakeSentry struct {
	mock.Mock
}

func (f *FakeSentry) CaptureMessage(message string) {
	f.Called(message)
}

func (f *FakeSentry) CaptureException(exception error) {
	f.Called(exception)
}

func (f *FakeSentry) Flush(timeout time.Duration) {
	f.Called(timeout)
}

func (f *FakeSentry) StartSpan(_ context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	args := f.Called(operation, options)

	return args.Get(0).(*sentry.Span)
}
