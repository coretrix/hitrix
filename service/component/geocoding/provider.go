package geocoding

import (
	"context"
	"googlemaps.github.io/maps"
)

type Provider interface {
	SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error)
	Geocode(ctx context.Context, address string, language string) (*Address, interface{}, error)
	ReverseGeocode(ctx context.Context, latLng *LatLng, language string) (*Address, interface{}, error)
	GetName() string
}
