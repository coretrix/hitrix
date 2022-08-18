# Sentry service
This service allow you to use sentry for logging events and performance tracking.

Register the service into your `main.go` file:
```go
registry.ServiceProviderSentry(tracesSampleRate *float64)
```

Access the service:
```go
service.DI().Sentry()
```

The methods that this service provide are:
```go
type ISentry interface {
    CaptureMessage(message string)
    Flush(timeout time.Duration)
    StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) *sentry.Span
}
```
You should call `CaptureMessage` when you want to send event to sentry

You should call `Flush` in your main file with defer, with 2 second timeout

You should call `StartSpan` when you want to start performance monitor
