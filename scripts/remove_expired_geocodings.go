package scripts

import (
	"context"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type RemoveExpiredGeocodingsScript struct {
}

func (script *RemoveExpiredGeocodingsScript) Run(_ context.Context, _ app.IExit) {
	now := service.DI().Clock().Now()

	fiveAM := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())

	if now.After(fiveAM) {
		return
	}

	ormService := service.DI().OrmEngine().Clone()

	where := beeorm.NewWhere("ExpiresAt < ?", now)

	geocodingEntities := make([]*entity.GeocodingEntity, 0)
	ormService.Search(where, beeorm.NewPager(1, 10000), &geocodingEntities)

	flusher := ormService.NewFlusher()

	for _, geocodingEntity := range geocodingEntities {
		flusher.Delete(geocodingEntity)
	}

	flusher.Flush()

	where = beeorm.NewWhere("ExpiresAt < ?", now)

	reverseGeocodingEntities := make([]*entity.GeocodingReverseEntity, 0)
	ormService.Search(where, beeorm.NewPager(1, 10000), &reverseGeocodingEntities)

	flusher = ormService.NewFlusher()

	for _, reverseGeocodingEntity := range reverseGeocodingEntities {
		flusher.Delete(reverseGeocodingEntity)
	}

	flusher.Flush()
}

func (script *RemoveExpiredGeocodingsScript) Interval() time.Duration {
	return time.Hour * 3
}

func (script *RemoveExpiredGeocodingsScript) Unique() bool {
	return true
}

func (script *RemoveExpiredGeocodingsScript) Description() string {
	return "remove expired geocodings"
}
