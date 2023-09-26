package scripts

import (
	"context"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type ClearExpiredGeocodingCache struct {
}

func (script *ClearExpiredGeocodingCache) Run(_ context.Context, ormService *beeorm.Engine, _ app.IExit) {
	now := service.DI().Clock().Now()

	fiveAM := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())

	if now.After(fiveAM) {
		//TODO sleep few hours
		return
	}

	where := beeorm.NewWhere("ExpiresAt < ?", now)

	geocodingEntities := make([]*entity.GeocodingCacheEntity, 0)
	ormService.Search(where, beeorm.NewPager(1, 10000), &geocodingEntities)

	flusher := ormService.NewFlusher()

	for _, geocodingEntity := range geocodingEntities {
		flusher.Delete(geocodingEntity)
	}

	flusher.Flush()

	where = beeorm.NewWhere("ExpiresAt < ?", now)

	reverseGeocodingEntities := make([]*entity.GeocodingReverseCacheEntity, 0)
	ormService.Search(where, beeorm.NewPager(1, 10000), &reverseGeocodingEntities)

	flusher = ormService.NewFlusher()

	for _, ReverseGeocodingCacheEntity := range reverseGeocodingEntities {
		flusher.Delete(ReverseGeocodingCacheEntity)
	}

	flusher.Flush()
}

func (script *ClearExpiredGeocodingCache) Interval() time.Duration {
	return time.Hour * 3
}

func (script *ClearExpiredGeocodingCache) Unique() bool {
	return true
}

func (script *ClearExpiredGeocodingCache) Description() string {
	return "clear expired geocoding cache"
}
