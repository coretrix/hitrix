# ClockWork
Provides you information about all queries executed by ORM in the current request. Also gives you ability to set your own log data.

It requires installation of Clockwork Chrome extension [ClockWork extension](https://chrome.google.com/webstore/detail/clockwork/dmggabnehkmmfmdffgajcflpdjlnoemp)

You need to add special header to activate this feature. My recommendation is to install also this extension [ModHeader extension](https://chrome.google.com/webstore/detail/modheader/idgpnmonknjnojddfkpgkljpfnnfcklj)
You should set header `CoreTrix` with value equal to the password you will set bellow in your yaml file

```yaml
clockwork:
    password: "your password here"

```

Register the service into your `main.go` file as context service:
```go 
registry.ServiceProviderClockWorkForContext()
```

Access the service:
```go
service.DI().ClockWorkForContext(ctx).GetLoggerDataSource().LogDebugString("key", "test")
```

There are 2 steps that also needs to be done:
1. To add this middleware

```go
	hitrixMiddleware "github.com/coretrix/hitrix/pkg/middleware"
	
	...
	
	hitrixMiddleware.Clockwork(ginEngine)
```

2. To add special route
```go
	hitrixController "github.com/coretrix/hitrix/pkg/controller"

    ...

	var clockwork *hitrixController.ClockworkController
	ginEngine.GET("/__clockwork/:id", clockwork.GetIndexAction)
```