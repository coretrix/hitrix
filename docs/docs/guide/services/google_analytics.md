# Google Analytics

This service is used for querying from Google Analytics

Register the service into your `main.go`. You need to provide function that init the analytics type (UA or GA4)

```go
    registry.ServiceProviderGoogleAnalytics(googleanalytics.NewGA4)
```

Then you need to download client configuration (credentials file) from your panel and put it in configs folder.
After that, you should put your ID of your GA property and config file name in config yaml file

Example:
```yml
google_analytics:
  config_file_name: Name-s0m3r4nd0mv4lu3.json
  property_id: 123456789
```


Access the service:
```go
service.DI().GoogleAnalytics().GetProvider(googleanalytics.GA4)
```
