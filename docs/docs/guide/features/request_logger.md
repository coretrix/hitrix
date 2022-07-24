# Request logger

Request logger `middleware` allows you to log every request that comes to the `API`

This is really helpful to debug your application

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

We created a `Cleaner` that will remove all rows in `request_logger` table older than 1 month. This will prevent your database to be fulfilled with logs

To enable it please put this code into your `single-instance-cron`

```go
    b := &hitrix.BackgroundProcessor{Server: s}
    b.RunAsyncRequestLoggerCleaner()
```

Using our Dev Panel you will be able easily to see all requests and search trough them