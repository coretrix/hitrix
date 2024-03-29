package featureflag

import (
	"sync"
	"time"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type cacheEntity struct {
	featureFlagEntity *entity.FeatureFlagEntity
	time              time.Time
}
type serviceFeatureFlagWithCache struct {
	sync.Mutex
	featureFlagService ServiceFeatureFlagInterface
	clockService       clock.IClock
	cache              map[string]*cacheEntity
}

func NewFeatureFlagWithCacheService(errorLoggerService errorlogger.ErrorLogger, clockService clock.IClock) ServiceFeatureFlagInterface {
	featureFlagService := NewFeatureFlagService(errorLoggerService)

	cachedService := &serviceFeatureFlagWithCache{
		featureFlagService: featureFlagService,
		clockService:       clockService,
		cache:              make(map[string]*cacheEntity),
	}

	return cachedService
}

func (s *serviceFeatureFlagWithCache) IsActive(ormService *datalayer.ORM, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	s.Lock()
	defer s.Unlock()

	if cacheEntry, has := s.cache[name]; has {
		if s.clockService.Now().Sub(cacheEntry.time) <= 5*time.Second {
			return cacheEntry.featureFlagEntity.Enabled && cacheEntry.featureFlagEntity.Registered
		}
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}

	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return false
	}

	s.cache[name] = &cacheEntity{
		featureFlagEntity: featureFlagEntity,
		time:              s.clockService.Now(),
	}

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlagWithCache) FailIfIsNotActive(ormService *datalayer.ORM, name string) error {
	return s.featureFlagService.FailIfIsNotActive(ormService, name)
}

func (s *serviceFeatureFlagWithCache) Enable(ormService *datalayer.ORM, name string) error {
	err := s.featureFlagService.Enable(ormService, name)
	s.Lock()
	delete(s.cache, name)
	s.Unlock()

	return err
}

func (s *serviceFeatureFlagWithCache) Disable(ormService *datalayer.ORM, name string) error {
	err := s.featureFlagService.Disable(ormService, name)
	s.Lock()
	delete(s.cache, name)
	s.Unlock()

	return err
}

func (s *serviceFeatureFlagWithCache) GetScriptsSingleInstance(ormService *datalayer.ORM) []app.IScript {
	return s.featureFlagService.GetScriptsSingleInstance(ormService)
}

func (s *serviceFeatureFlagWithCache) GetScriptsMultiInstance(ormService *datalayer.ORM) []app.IScript {
	return s.featureFlagService.GetScriptsMultiInstance(ormService)
}

func (s *serviceFeatureFlagWithCache) Register(featureFlags ...IFeatureFlag) {
	s.featureFlagService.Register(featureFlags...)
}

func (s *serviceFeatureFlagWithCache) Sync(ormService *datalayer.ORM, clockService clock.IClock) {
	s.featureFlagService.Sync(ormService, clockService)
}
