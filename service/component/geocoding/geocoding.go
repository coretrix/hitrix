package geocoding

import (
	"context"
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive, but just fine for caching
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	Geocode(ctx context.Context, ormService *datalayer.ORM, address string, language Language) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *datalayer.ORM, latLng *LatLng, language Language) (*Address, error)
	CutCoordinates(float float64, precision int) (float64, error)
}

type Address struct {
	Found     bool
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
	useCaching      bool
	cacheTTLMinDays int
	cacheTTLMaxDays int
	clock           clock.IClock
	provider        Provider
}

func NewGeocoding(
	useCaching bool,
	cacheTTLMinDays int,
	cacheTTLMaxDays int,
	clock clock.IClock,
	provider Provider,
) IGeocoding {
	return &Geocoding{
		useCaching:      useCaching,
		cacheTTLMinDays: cacheTTLMinDays,
		cacheTTLMaxDays: cacheTTLMaxDays,
		clock:           clock,
		provider:        provider,
	}
}

func (g *Geocoding) Geocode(ctx context.Context, ormService *datalayer.ORM, address string, language Language) (*Address, error) {
	languageEnum, ok := languageToEnumMapping[language]
	if !ok {
		return nil, fmt.Errorf("language %s not supported", language)
	}

	address = strings.TrimSpace(address)

	if g.useCaching {
		geocodingEntity := &entity.GeocodingCacheEntity{}
		if ormService.CachedSearchOne(geocodingEntity, "CachedQueryAddressHashLanguage", g.getAddressHash(address), language) {
			return &Address{
				Found:     true,
				FromCache: true,
				Address:   address,
				Language:  Language(geocodingEntity.Language),
				Location: &LatLng{
					Lat: geocodingEntity.Lat,
					Lng: geocodingEntity.Lng,
				},
			}, nil
		}
	}

	geocodedAddress, providerRawResponse, err := g.provider.Geocode(ctx, address, language)
	if err != nil {
		return nil, err
	}

	if g.useCaching && geocodedAddress.Found {
		now := g.clock.Now()

		ormService.Flush(&entity.GeocodingCacheEntity{
			Lat:         geocodedAddress.Location.Lat,
			Lng:         geocodedAddress.Location.Lng,
			Address:     address,
			AddressHash: g.getAddressHash(address),
			Language:    languageEnum,
			Provider:    g.provider.GetName(),
			RawResponse: providerRawResponse,
			ExpiresAt:   now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return geocodedAddress, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService *datalayer.ORM, latLng *LatLng, language Language) (*Address, error) {
	languageEnum, ok := languageToEnumMapping[language]
	if !ok {
		return nil, fmt.Errorf("language %s not supported", language)
	}

	cacheLat := latLng.Lat
	cacheLng := latLng.Lng

	if g.useCaching {
		var err error

		cacheLat, err = g.CutCoordinates(cacheLat, 5)
		if err != nil {
			return nil, err
		}

		cacheLng, err = g.CutCoordinates(cacheLng, 5)
		if err != nil {
			return nil, err
		}
	}

	if g.useCaching {
		ReverseGeocodingCacheEntity := &entity.GeocodingReverseCacheEntity{}
		if ormService.CachedSearchOne(ReverseGeocodingCacheEntity, "CachedQueryLatLngLanguage", cacheLat, cacheLng, language) {
			return &Address{
				Found:     true,
				FromCache: true,
				Address:   ReverseGeocodingCacheEntity.Address,
				Language:  Language(ReverseGeocodingCacheEntity.Language),
				Location: &LatLng{
					Lat: ReverseGeocodingCacheEntity.Lat,
					Lng: ReverseGeocodingCacheEntity.Lng,
				},
			}, nil
		}
	}

	geocodedAddress, providerRawResponse, err := g.provider.ReverseGeocode(ctx, latLng, language)
	if err != nil {
		return nil, err
	}

	if g.useCaching && geocodedAddress.Found {
		now := g.clock.Now()

		ormService.Flush(&entity.GeocodingReverseCacheEntity{
			Lat:         cacheLat,
			Lng:         cacheLng,
			Address:     strings.TrimSpace(geocodedAddress.Address),
			Language:    languageEnum,
			Provider:    g.provider.GetName(),
			RawResponse: providerRawResponse,
			ExpiresAt:   now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:   now,
		})
	}

	return geocodedAddress, nil
}

func (g *Geocoding) CutCoordinates(float float64, precision int) (float64, error) {
	asString := fmt.Sprintf("%.8f", float)
	asStringParts := strings.Split(asString, ".")

	return strconv.ParseFloat(asString[0:len(asStringParts[0])+1+precision], 64)
}

func (g *Geocoding) getCacheTTL(cacheTTLMinDays, cacheTTLMaxDays int) int {
	intVal := big.NewInt(int64(cacheTTLMaxDays - cacheTTLMinDays))

	randomInt, err := rand.Int(rand.Reader, intVal)
	if err != nil {
		panic(err)
	}

	return int(randomInt.Int64() + int64(cacheTTLMinDays))
}

func (g *Geocoding) getAddressHash(address string) string {
	//nolint // G401: Use of weak cryptographic primitive , but just fine for caching
	hash := md5.Sum([]byte(address))

	return hex.EncodeToString(hash[:])
}
