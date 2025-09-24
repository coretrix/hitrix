package geocoding

import (
	"context"
	"errors"
	"googlemaps.github.io/maps"

	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive, but just fine for caching
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error)
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language string) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language string) (*Address, error)
	CutCoordinates(float float64, precision int) (float64, error)
}

type Address struct {
	Found                    bool
	FromCache                bool
	AdministrativeAreaLevel1 string
	CityName                 string
	Address                  string
	Language                 string
	Location                 *LatLng
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

func (g *Geocoding) SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error) {
	return g.provider.SnapToRoad(ctx, dto)
}

func (g *Geocoding) Geocode(ctx context.Context, ormService *beeorm.Engine, address string, language string) (*Address, error) {
	address = strings.TrimSpace(address)

	if g.useCaching {
		geocodingEntity := &entity.GeocodingCacheEntity{}
		if ormService.CachedSearchOne(geocodingEntity, "CachedQueryAddressHashLanguage", g.getAddressHash(address), language) {
			return &Address{
				Found:                    true,
				FromCache:                true,
				AdministrativeAreaLevel1: geocodingEntity.AdministrativeAreaLevel1,
				CityName:                 geocodingEntity.CityName,
				Address:                  address,
				Language:                 geocodingEntity.Language,
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

		geocodingCacheEntity := &entity.GeocodingCacheEntity{
			Lat:         geocodedAddress.Location.Lat,
			Lng:         geocodedAddress.Location.Lng,
			Address:     address,
			AddressHash: g.getAddressHash(address),
			Language:    language,
			Provider:    g.provider.GetName(),
			RawResponse: providerRawResponse,
			ExpiresAt:   now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:   now,
		}

		administrativeAreaL1, cityName := g.extractRegionAndCity(providerRawResponse.(maps.GeocodingResult))
		geocodingCacheEntity.AdministrativeAreaLevel1 = administrativeAreaL1
		geocodingCacheEntity.CityName = cityName

		ormService.Flush(geocodingCacheEntity)
	}

	return geocodedAddress, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng, language string) (*Address, error) {
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

		reverseGeocodingCacheEntity := &entity.GeocodingReverseCacheEntity{}

		found := ormService.CachedSearchOne(reverseGeocodingCacheEntity, "CachedQueryLatLngLanguage", cacheLat, cacheLng, language)
		if found {
			return &Address{
				Found:                    true,
				FromCache:                true,
				AdministrativeAreaLevel1: reverseGeocodingCacheEntity.AdministrativeAreaLevel1,
				CityName:                 reverseGeocodingCacheEntity.CityName,
				Address:                  reverseGeocodingCacheEntity.Address,
				Language:                 reverseGeocodingCacheEntity.Language,
				Location: &LatLng{
					Lat: reverseGeocodingCacheEntity.Lat,
					Lng: reverseGeocodingCacheEntity.Lng,
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

		geocodingReverseCacheEntity := &entity.GeocodingReverseCacheEntity{
			Lat:         cacheLat,
			Lng:         cacheLng,
			Address:     strings.TrimSpace(geocodedAddress.Address),
			Language:    language,
			Provider:    g.provider.GetName(),
			RawResponse: providerRawResponse,
			ExpiresAt:   now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:   now,
		}

		administrativeAreaL1, cityName := g.extractRegionAndCity(providerRawResponse.(maps.GeocodingResult))
		geocodingReverseCacheEntity.AdministrativeAreaLevel1 = administrativeAreaL1
		geocodingReverseCacheEntity.CityName = cityName

		err := ormService.FlushWithCheck()

		//TODO Krasi: needed due to issue when localCache and redisCache used together
		if err != nil {
			var duplicateKeyError *beeorm.DuplicatedKeyError
			if errors.As(err, &duplicateKeyError) {
				if duplicateKeyError.Index != "geocoding_reverse_cache.Lat_Lng_Language" {
					panic(err)
				}

				geocodedAddress.Found = true

				return geocodedAddress, nil
			}

			panic(err)
		}
	}

	return geocodedAddress, nil
}

func (g *Geocoding) extractRegionAndCity(result maps.GeocodingResult) (region, city string) {
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

	return
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
