package geocoding

import (
	"context"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language string) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language string) (*Address, error)
}

type Address struct {
	FromCache bool
	Address   string
	Language  string
	Location  *LatLng
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Geocoding struct {
	useCaching      bool
	cacheExpiryDays int64
	clock           clock.IClock
	provider        Provider
}

func NewGeocoding(useCaching bool, cacheExpiryDays int64, clock clock.IClock, provider Provider) IGeocoding {
	return &Geocoding{useCaching: useCaching, cacheExpiryDays: cacheExpiryDays, clock: clock, provider: provider}
}

func (g *Geocoding) Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language string) (*Address, error) {
	if g.useCaching {
		geocodingEntity := &entity.GeocodingEntity{}
		if ormService.CachedSearchOne(geocodingEntity, "CachedQueryAddressLanguage", address, language) {
			return &Address{
				FromCache: true,
				Address:   geocodingEntity.Address,
				Language:  language,
				Location: &LatLng{
					Lat: geocodingEntity.Lat,
					Lng: geocodingEntity.Lng,
				},
			}, nil
		}
	}

	addressResult, rawResponse, err := g.provider.Geocode(ctx, language, address)
	if err != nil {
		return nil, err
	}

	now := g.clock.Now()

	if g.useCaching {
		ormService.Flush(&entity.GeocodingEntity{
			Lat:         addressResult.Location.Lat,
			Lng:         addressResult.Location.Lng,
			Address:     addressResult.Address,
			Language:    language,
			Provider:    g.provider.GetName(),
			RawResponse: rawResponse,
			ExpiresAt:   now.Add(time.Duration(g.cacheExpiryDays) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return addressResult, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language string) (*Address, error) {
	if g.useCaching {
		reverseGeocodingEntity := &entity.ReverseGeocodingEntity{}
		if ormService.CachedSearchOne(reverseGeocodingEntity, "CachedQueryLatLngLanguage", latLng.Lat, latLng.Lng, language) {
			return &Address{
				FromCache: true,
				Address:   reverseGeocodingEntity.Address,
				Language:  language,
				Location: &LatLng{
					Lat: reverseGeocodingEntity.Lat,
					Lng: reverseGeocodingEntity.Lng,
				},
			}, nil
		}
	}

	addressResult, rawResponse, err := g.provider.ReverseGeocode(ctx, latLng, language)
	if err != nil {
		return nil, err
	}

	now := g.clock.Now()

	if g.useCaching {
		ormService.Flush(&entity.ReverseGeocodingEntity{
			Lat:         addressResult.Location.Lat,
			Lng:         addressResult.Location.Lng,
			Address:     addressResult.Address,
			Language:    language,
			Provider:    g.provider.GetName(),
			RawResponse: rawResponse,
			ExpiresAt:   now.Add(time.Duration(g.cacheExpiryDays) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return addressResult, nil
}
