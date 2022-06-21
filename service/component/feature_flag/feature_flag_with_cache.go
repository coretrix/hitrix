package featureflag

import (
	"sync"
	"time"

	"github.com/latolukasz/beeorm"

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

func (s *serviceFeatureFlagWithCache) IsActive(ormService *beeorm.Engine, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	if cacheEntry, has := s.cache[name]; has {
		if s.clockService.Now().Sub(cacheEntry.time) <= 5*time.Second {
			return cacheEntry.featureFlagEntity.Enabled && cacheEntry.featureFlagEntity.Registered
		}
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return false
	}
	s.Lock()
	s.cache[name] = &cacheEntity{
		featureFlagEntity: featureFlagEntity,
		time:              s.clockService.Now(),
	}
	defer s.Unlock()

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlagWithCache) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	return s.featureFlagService.FailIfIsNotActive(ormService, name)
}

func (s *serviceFeatureFlagWithCache) Enable(ormService *beeorm.Engine, name string) error {
	return s.featureFlagService.Enable(ormService, name)
}

func (s *serviceFeatureFlagWithCache) Disable(ormService *beeorm.Engine, name string) error {
	return s.featureFlagService.Disable(ormService, name)
}

func (s *serviceFeatureFlagWithCache) GetScriptsSingleInstance(ormService *beeorm.Engine) []app.IScript {
	return s.featureFlagService.GetScriptsSingleInstance(ormService)
}

func (s *serviceFeatureFlagWithCache) GetScriptsMultiInstance(ormService *beeorm.Engine) []app.IScript {
	return s.featureFlagService.GetScriptsMultiInstance(ormService)
}

func (s *serviceFeatureFlagWithCache) Register(featureFlags ...IFeatureFlag) {
	s.featureFlagService.Register(featureFlags...)
}

func (s *serviceFeatureFlagWithCache) Sync(ormService *beeorm.Engine, clockService clock.IClock) {
	s.featureFlagService.Sync(ormService, clockService)
}
