package geocoding

import (
	"context"
)

type Provider interface {
	Geocode(ctx context.Context, address string, language Language) (*Address, interface{}, error)
	ReverseGeocode(ctx context.Context, latLng *LatLng, language Language) (*Address, interface{}, error)
	GetName() string
}
