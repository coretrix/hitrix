# Geocoding

This service is used for geocoding and reverse geocoding. It supports multiple providers. For now only
Google Maps provider is implemented.


You should put your api key from Google Maps in config:

```yml
geocoding:
  use_caching: true
  cache_expiry_days: 10
  cache_lat_lng_floating_point_precision: 2
  google_maps:
    api_key: some_key
```

Note 1: if you decide to use caching functionality, you need to run script `RemoveExpiredGeocodingsScript`
in your project. This script will delete expired cache after days setting `cache_expiry_days`.

Note 2: If you include `cache_lat_lng_floating_point_precision` in your config, when caching
coordinates the service will truncate the lat and lng values to this setting.

Access the service:
```go
service.DI().Geocoding()
```

The service exposes 2 methods that you can use:

```go
type IGeocoding interface {
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language Language) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language Language) (*Address, error)
}
```
