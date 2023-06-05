package geocoding

import (
	"context"
	"fmt"
	"github.com/coretrix/hitrix/pkg/helper"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language Language) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language Language) (*Address, error)
}

type Address struct {
	FromCache bool
	Address   string
	Language  Language
	Location  *LatLng
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Geocoding struct {
	useCaching                           bool
	cacheExpiryDays                      int64
	cacheLatLngFloatingPointPrecision    int
	useCacheLatLngFloatingPointPrecision bool
	clock                                clock.IClock
	provider                             Provider
}

func NewGeocoding(
	useCaching bool,
	cacheExpiryDays int64,
	cacheLatLngFloatingPointPrecision int,
	useCacheLatLngFloatingPointPrecision bool,
	clock clock.IClock,
	provider Provider,
) IGeocoding {
	return &Geocoding{
		useCaching:                           useCaching,
		cacheExpiryDays:                      cacheExpiryDays,
		cacheLatLngFloatingPointPrecision:    cacheLatLngFloatingPointPrecision,
		useCacheLatLngFloatingPointPrecision: useCacheLatLngFloatingPointPrecision,
		clock:                                clock,
		provider:                             provider,
	}
}

func (g *Geocoding) Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language Language) (*Address, error) {
	languageEnum, ok := languageToEnumMapping[language]
	if !ok {
		return nil, fmt.Errorf("language %s not supported", language)
	}

	if g.useCaching {
		geocodingEntity := &entity.GeocodingEntity{}
		if ormService.CachedSearchOne(geocodingEntity, "CachedQueryAddressLanguage", address, language) {
			return &Address{
				FromCache: true,
				Address:   geocodingEntity.Address,
				Language:  Language(geocodingEntity.Language),
				Location: &LatLng{
					Lat: geocodingEntity.Lat,
					Lng: geocodingEntity.Lng,
				},
			}, nil
		}
	}

	addressResult, rawResponse, err := g.provider.Geocode(ctx, address, language)
	if err != nil {
		return nil, err
	}

	now := g.clock.Now()

	if g.useCaching {
		ormService.Flush(&entity.GeocodingEntity{
			Lat:         addressResult.Location.Lat,
			Lng:         addressResult.Location.Lng,
			Address:     addressResult.Address,
			Language:    languageEnum,
			Provider:    g.provider.GetName(),
			RawResponse: rawResponse,
			ExpiresAt:   now.Add(time.Duration(g.cacheExpiryDays) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return addressResult, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language Language) (*Address, error) {
	languageEnum, ok := languageToEnumMapping[language]
	if !ok {
		return nil, fmt.Errorf("language %s not supported", language)
	}

	cacheLat := latLng.Lat
	cacheLng := latLng.Lng

	if g.useCacheLatLngFloatingPointPrecision {
		cacheLat = helper.ToFixed(cacheLat, g.cacheLatLngFloatingPointPrecision)
		cacheLng = helper.ToFixed(cacheLng, g.cacheLatLngFloatingPointPrecision)
	}

	if g.useCaching {
		reverseGeocodingEntity := &entity.ReverseGeocodingEntity{}
		if ormService.CachedSearchOne(reverseGeocodingEntity, "CachedQueryLatLngLanguage", cacheLat, cacheLng, language) {
			return &Address{
				FromCache: true,
				Address:   reverseGeocodingEntity.Address,
				Language:  Language(reverseGeocodingEntity.Language),
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
			Lat:         cacheLat,
			Lng:         cacheLng,
			Address:     addressResult.Address,
			Language:    languageEnum,
			Provider:    g.provider.GetName(),
			RawResponse: rawResponse,
			ExpiresAt:   now.Add(time.Duration(g.cacheExpiryDays) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return addressResult, nil
}