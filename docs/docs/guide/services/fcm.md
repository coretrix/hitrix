# Firebase cloud messaging (FCM) service
This service is used for sending different types of push notifications

Register the service into your `main.go` file:
```go
registry.ServiceProviderFCM(),
```

Config sample:

expose `FIREBASE_CONFIG="path/to/service-account-file.json"`
