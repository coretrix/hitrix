package featureflag

import (
	"sync"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type cache struct {
	ttl          int64
	clockService clock.IClock
	cacheMap     sync.Map
}

type cacheEntry struct {
	featureFlagEntity *entity.FeatureFlagEntity
	time              time.Time
}

func (c *cache) Store(name string, value cacheEntry) {
	c.cacheMap.Store(name, value)
}
func (c *cache) Load(name string) (*entity.FeatureFlagEntity, bool) {
	value, has := c.cacheMap.Load(name)
	if !has {
		return nil, false
	}

	cacheEntry := value.(cacheEntry)
	if c.clockService.Now().Sub(cacheEntry.time) <= time.Second*time.Duration(c.ttl) {
		return cacheEntry.featureFlagEntity, true
	}

	return nil, false
}

func (c *cache) Delete(name string) {
	c.cacheMap.Delete(name)
}
