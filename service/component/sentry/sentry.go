package sentry

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

type ISentry interface {
	CaptureMessage(message string)
	Flush(timeout time.Duration)
	StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span
}

type v struct {
}

func Init(dsn, release string, tracesSampleRate *float64) ISentry {
	tracesSampleRateValue := 0.0
	if tracesSampleRate != nil {
		tracesSampleRateValue = *tracesSampleRate
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: tracesSampleRateValue,
		Release:          release,
	})
	if err != nil {
		panic(err)
	}

	return &v{}
}

func (v *v) CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}

func (v *v) Flush(timeout time.Duration) {
	sentry.Flush(timeout)
}

func (v *v) StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	return sentry.StartSpan(ctx, operation, options...)
}
