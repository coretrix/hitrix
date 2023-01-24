# Geocoding

This service is used for geocoding and reverse geocoding. It supports multiple providers. For now only
Google Maps provider is implemented.


You should put your api key from Google Maps in config:

```yml
geocoding:
  google_maps:
    api_key: some_key
```

Access the service:
```go
service.DI().Geocoding()
```
