# API logger service
This service is used to track every api request and response.
For example it can be used in any other service as SMS service, Stripe service and so on. Using it you will have a history of all requests and responses
and it will help you even in case you need to debug something.

Register the service into your `main.go` file:
```go
registry.APILogger(&entity.APILogEntity{}),
```

Access the service:
```go
service.DI().APILogger()
```

All the data it will be saved into the `APILogEntity` entity.

The methods that this service provide are:
```go
type APILogger interface {
	LogStart(logType string, request interface{})
	LogError(message string, response interface{})
	LogSuccess(response interface{})
}
```
You should call `LogStart` before you send request to the api

You should call `LogError` in case api return you error

You should call `LogSuccess` in case api return you success
