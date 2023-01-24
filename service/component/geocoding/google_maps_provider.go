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

func (g *GoogleMapsProvider) Geocode(ctx context.Context, address string, language string) (*Address, interface{}, error) {
	req := &maps.GeocodingRequest{
		Language: language,
		Address:  address,
	}

	results, err := g.client.Geocode(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	if len(results) == 0 {
		return &Address{
			Address:  address,
			Language: language,
			Location: &LatLng{
				Lat: 0,
				Lng: 0,
			},
		}, results, err
	}

	return &Address{
		Address:  results[0].FormattedAddress,
		Language: language,
		Location: &LatLng{
			Lat: results[0].Geometry.Location.Lat,
			Lng: results[0].Geometry.Location.Lng,
		},
	}, results, err
}

func (g *GoogleMapsProvider) ReverseGeocode(ctx context.Context, latLng *LatLng, language string) (*Address, interface{}, error) {
	req := &maps.GeocodingRequest{
		Language: language,
		LatLng: &maps.LatLng{
			Lat: latLng.Lat,
			Lng: latLng.Lng,
		},
	}

	results, err := g.client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	if len(results) == 0 {
		return &Address{
			Address:  "",
			Language: language,
			Location: &LatLng{
				Lat: latLng.Lat,
				Lng: latLng.Lng,
			},
		}, results, err
	}

	return &Address{
		Address:  results[0].FormattedAddress,
		Language: language,
		Location: &LatLng{
			Lat: results[0].Geometry.Location.Lat,
			Lng: results[0].Geometry.Location.Lng,
		},
	}, results, err
}

func (g *GoogleMapsProvider) GetName() string {
	return "google maps"
}
