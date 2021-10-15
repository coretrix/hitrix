# Error Logger
Used to save unhandled errors in error log. Hitrix use `recovery` function to handle those errors.
If you setup Slack service you also going to receive notifications in your slack

Register the service into your `main.go` file:
```go 
registry.ServiceProviderErrorLogger()
```

Access the service:
```go
service.DI().ErrorLogger()
```

It can be used to save custom errors as well:
```go
        errorLoggerService := ioc.GetErrorLoggerService()
		
		errorLoggerService.LogErrorWithRequest(c, err) //if you provide context we will save request body as well
		errorLoggerService.LogError(err)
```