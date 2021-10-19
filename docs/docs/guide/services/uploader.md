# Uploader

This service uses TUS protocol to enable fast resumable and multi-part upload of big files.
It provides an easy interface for plug-in whatever data store and locker you want to implement.
Currently, Amazon S3 data store and Redis locker are implemented. For Amazon data store to work,
you need to register Amazon S3 service before this one, also for Redis locker to work, you need
to register orm service background before this one.

Register the service into your `main.go` file:
```go
registry.ServiceProviderUploader(tusd.Config{...}, datastore.GetAmazonS3Store, locker.GetRedisLocker)
```

Access the service:
```go
service.DI().Uploader()
```

Hitrix also provides REST uploader controller which you can register all handler methods in your
router:

```go
var uploaderController *hitrixController.UploaderController
uploaderGroup := ginEngine.Group("/files/")
uploaderGroup.Use(middleware.AuthorizeWithHeaderStrict())
{
	uploaderGroup.POST("", uploaderController.PostFileAction)
	uploaderGroup.HEAD(":id", uploaderController.HeadFile)
	uploaderGroup.PATCH(":id", uploaderController.PatchFile)
	uploaderGroup.GET(":id", uploaderController.GetFileAction)
	uploaderGroup.DELETE(":id", uploaderController.DeleteFile)
}
```

Also you need bucket name in config:

````yml
uploader:
  bucket: media
````