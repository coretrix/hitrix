package geocoding

import (
	"context"
)

type Provider interface {
	Geocode(ctx context.Context, address string) ([]*Address, error)
	ReverseGeocode(ctx context.Context, latLng *LatLng) ([]*Address, error)
}
