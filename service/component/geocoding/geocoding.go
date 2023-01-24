package geocoding

import (
	"context"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	Geocode(ctx context.Context, ormService *beeorm.Engine, address string) ([]*Address, error)
	ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng) ([]*Address, error)
}

type Address struct {
	Address  string
	Location *LatLng
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Geocoding struct {
	useCaching bool
	clock      clock.IClock
	provider   Provider
}

func NewGeocoding(useCaching bool, clock clock.IClock, provider Provider) IGeocoding {
	return &Geocoding{useCaching: useCaching, clock: clock, provider: provider}
}

func (g *Geocoding) Geocode(ctx context.Context, ormService *beeorm.Engine, address string) ([]*Address, error) {
	if g.useCaching {
		q := beeorm.NewRedisSearchQuery()
		q.FilterString("Address", address)
		q.Sort("ID", true)

		geocodingEntities := make([]*entity.GeocodingEntity, 0)
		ormService.RedisSearch(&geocodingEntities, q, beeorm.NewPager(1, 1000))

		if len(geocodingEntities) != 0 {
			return adaptEntitiesToAddresses(geocodingEntities), nil
		}
	}

	addresses, err := g.provider.Geocode(ctx, address)
	if err != nil {
		return nil, err
	}

	now := g.clock.Now()

	if g.useCaching && len(addresses) != 0 {
		flusher := ormService.NewFlusher()

		for _, addressResult := range addresses {
			flusher.Track(&entity.GeocodingEntity{
				Lat:       addressResult.Location.Lat,
				Lng:       addressResult.Location.Lng,
				Address:   addressResult.Address,
				CreatedAt: now,
			})
		}

		flusher.Flush()
	}

	return addresses, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService *beeorm.Engine, latLng *LatLng) ([]*Address, error) {
	if g.useCaching {
		q := beeorm.NewRedisSearchQuery()
		q.FilterFloat("Lat", latLng.Lat)
		q.FilterFloat("Lng", latLng.Lng)
		q.Sort("ID", true)

		geocodingEntities := make([]*entity.GeocodingEntity, 0)
		ormService.RedisSearch(&geocodingEntities, q, beeorm.NewPager(1, 1000))

		if len(geocodingEntities) != 0 {
			return adaptEntitiesToAddresses(geocodingEntities), nil
		}
	}

	addresses, err := g.provider.ReverseGeocode(ctx, latLng)
	if err != nil {
		return nil, err
	}

	now := g.clock.Now()

	if g.useCaching && len(addresses) != 0 {
		flusher := ormService.NewFlusher()

		for _, addressResult := range addresses {
			flusher.Track(&entity.GeocodingEntity{
				Lat:       addressResult.Location.Lat,
				Lng:       addressResult.Location.Lng,
				Address:   addressResult.Address,
				CreatedAt: now,
			})
		}

		flusher.Flush()
	}

	return addresses, nil
}

func adaptEntitiesToAddresses(geocodingEntities []*entity.GeocodingEntity) []*Address {
	addresses := make([]*Address, len(geocodingEntities))

	for i, geocodingEntity := range geocodingEntities {
		addresses[i] = adaptEntityToAddress(geocodingEntity)
	}

	return addresses
}

func adaptEntityToAddress(geocodingEntity *entity.GeocodingEntity) *Address {
	return &Address{
		Address: geocodingEntity.Address,
		Location: &LatLng{
			Lat: geocodingEntity.Lat,
			Lng: geocodingEntity.Lng,
		},
	}
}
