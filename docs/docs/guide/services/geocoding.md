# Geocoding

This service is used for geocoding and reverse geocoding. It supports multiple providers. For now only
Google Maps provider is implemented.


You should put your api key from Google Maps in config:

```yml
geocoding:
  use_caching: true
  cache_ttl_min_days: 5
  cache_ttl_max_days: 10
  google_maps:
    api_key: some_key
```

Note 1: if you decide to use caching functionality, you need to run script `ClearExpiredGeocodingCache`
in your project. This script will delete expired cache.

Note 2: if you decide to use caching functionality, lat/lng are cut (not rounded) to 5 decimal place

Access the service:
```go
service.DI().Geocoding()
```

The service exposes 3 methods that you can use:

```go
type IGeocoding interface {
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language Language) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language Language) (*Address, error)
    CutCoordinates(float float64, precision int) (float64, error)
}
```
