package geocoding

import (
	"context"
)

type Provider interface {
	Geocode(ctx context.Context, address string, language string) (*Address, interface{}, error)
	ReverseGeocode(ctx context.Context, latLng *LatLng, language string) (*Address, interface{}, error)
	GetName() string
}
