# Localizer service
Localizer provides you a simple translation service that can pull and push translation pairs from local (file) and external sources (online services).

Currently localizer supports only [POEditor](https://poeditor.com) online source.

Localizer using a bucket key to separate and manage translation parts of your app.

First you need these in your app config:
```yaml
translation:
  poeditor:
    api_key: ENV[POEDITOR_API_KEY]
    project_id: ENV[POEDITOR_PROJECT_ID]
    language: ENV[POEDITOR_LANGUAGE]
```

Register the service into your `main.go` file:

```go
registry.ServiceProviderLocalizer()
```

Access the service:
```go
service.DI().LocalizerService()
```

Loading translation pairs from map:
```go
bucketKey := "greet-service"
append := false // append or replace?
pairs := map[string]string{
  "app_name": "My App Name",
  "loading_text": "Loading ...",
}
localizerService.LoadBucketFromMap(
  bucketKey, 
  pairs, 
  append,
)
```
Using `Localize()` function to translate a key:
```go
appName, err := localizerService.Localize(bucketKey, "app_name")
if err !nil {
  // handle error
}
```
Loading translation pairs from local file:
```go
localizerService.LoadBucketFromFile(
  bucketKey,
  "locales/greet.en.json",
  append,
)
```
Pull the translations from external source:
```go
err := localizerService.PullBucketFromSource(bucketKey, append)
if err != nil {
  log.Fatal(err)
}
```
Push translations to external source:
```go
err := localizerService.PushBucketToSource(bucketKey)
if err != nil {
  // handle error
}
```
