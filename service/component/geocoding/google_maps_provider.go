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
	response, err := g.client.Geocode(
		ctx,
		&maps.GeocodingRequest{
			Language: language,
			Address:  address,
		})
	if err != nil {
		return nil, nil, err
	}

	if len(response) == 0 {
		return &Address{
			Address:  address,
			Language: language,
			Location: &LatLng{
				Lat: 0,
				Lng: 0,
			},
		}, response, nil
	}

	administrativeAreaL1, cityName := g.extractRegionAndCity(response)

	return &Address{
		Found:                    true,
		AdministrativeAreaLevel1: administrativeAreaL1,
		CityName:                 cityName,
		Address:                  response[0].FormattedAddress,
		Language:                 language,
		Location: &LatLng{
			Lat: response[0].Geometry.Location.Lat,
			Lng: response[0].Geometry.Location.Lng,
		},
	}, response, nil
}

func (g *GoogleMapsProvider) extractRegionAndCity(results []maps.GeocodingResult) (region, city string) {
	for _, result := range results {
		for _, comp := range result.AddressComponents {
			for _, t := range comp.Types {
				switch t {
				case "administrative_area_level_1":
					region = comp.LongName
				case "locality":
					city = comp.LongName
				}
			}
		}
	}

	return
}

func (g *GoogleMapsProvider) ReverseGeocode(ctx context.Context, latLng *LatLng, language string) (*Address, interface{}, error) {
	response, err := g.client.ReverseGeocode(
		ctx,
		&maps.GeocodingRequest{
			Language: language,
			LatLng: &maps.LatLng{
				Lat: latLng.Lat,
				Lng: latLng.Lng,
			},
		})
	if err != nil {
		return nil, nil, err
	}

	if len(response) == 0 {
		return &Address{
			Language: language,
			Location: &LatLng{
				Lat: latLng.Lat,
				Lng: latLng.Lng,
			},
		}, response, nil
	}

	administrativeAreaL1, cityName := g.extractRegionAndCity(response)

	return &Address{
		Found:                    true,
		AdministrativeAreaLevel1: administrativeAreaL1,
		CityName:                 cityName,
		Address:                  response[0].FormattedAddress,
		Language:                 language,
		Location: &LatLng{
			Lat: response[0].Geometry.Location.Lat,
			Lng: response[0].Geometry.Location.Lng,
		},
	}, response, nil
}

func (g *GoogleMapsProvider) SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error) {
	snapResponse, err := g.client.SnapToRoad(ctx, dto)

	if err != nil {
		return nil, err
	}

	return snapResponse, nil
}

func (g *GoogleMapsProvider) GetName() string {
	return "google_maps"
}
