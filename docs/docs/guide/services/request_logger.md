# Request Logger service

This service is used for logging all upcoming and outgoing requests

Register the service into your `main.go` file:
```go
registry.ServiceProviderRequestLogger()
```
and you should register the entity `RequestLoggerEntity` into the ORM

Access the service:
```go
service.DI().RequestLogger()
```

The functions this service provide are:
```go
	LogRequest(ormService *datalayer.DataLayer, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
    LogResponse(ormService *datalayer.DataLayer, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
```
They can be used to log any outgoing requests you send

Also you are able to enable middleware which will log all incoming requests

```go
middleware.RequestLogger(ginEngine, func(context *gin.Context, requestEntity *entity.RequestLoggerEntity) {
			userService := ioc.GetUserService()
			session, hasSession := userService.GetSession(context.Request.Context())

			if hasSession && session.User != nil {
				requestEntity.UserID = session.User.ID
			}
		})
```

The second parameter (anonymous function) is called `extender` and it is used to save extra param to `request_logger` table like logged user id

If you want to use this `middleware` please do not forget to register the entity `RequestLoggerEntity`

We created a `Cleaner` that will remove all rows in `request_logger` table older than 30 days by default. This will prevent your database to be fulfilled with logs
If you want to change ttl to other value you can do it from your config file like on the example bellow:

```yaml
request_logger:
  ttl_in_days: 5
```
To enable it please put this code into your `single-instance-cron`

```go
    b := &hitrix.BackgroundProcessor{Server: s}
    b.RunAsyncRequestLoggerCleaner()
```

Using our Dev Panel you will be able easily to see all requests and search trough them
