package geocoding

import (
	"context"

	"googlemaps.github.io/maps"
)

type GoogleMapsProvider struct {
	client *maps.Client
}

func NewGoogleMapsProvider(apiKey string) Provider {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	return &GoogleMapsProvider{client: client}
}

func (g *GoogleMapsProvider) Geocode(ctx context.Context, address string) ([]*Address, error) {
	req := &maps.GeocodingRequest{
		Address: address,
	}

	results, err := g.client.Geocode(ctx, req)
	if err != nil {
		return nil, err
	}

	addresses := make([]*Address, len(results))

	for i, result := range results {
		addresses[i] = &Address{
			Address: result.FormattedAddress,
			Location: &LatLng{
				Lat: result.Geometry.Location.Lat,
				Lng: result.Geometry.Location.Lng,
			},
		}
	}

	return addresses, err
}

func (g *GoogleMapsProvider) ReverseGeocode(ctx context.Context, latLng *LatLng) ([]*Address, error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: latLng.Lat,
			Lng: latLng.Lng,
		},
	}

	results, err := g.client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, err
	}

	addresses := make([]*Address, len(results))

	for i, result := range results {
		addresses[i] = &Address{
			Address: result.FormattedAddress,
			Location: &LatLng{
				Lat: result.Geometry.Location.Lat,
				Lng: result.Geometry.Location.Lng,
			},
		}
	}

	return addresses, err
}
