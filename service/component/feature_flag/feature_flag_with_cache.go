package featureflag

import (
	"fmt"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type serviceFeatureFlagWithCache struct {
	featureFlagService ServiceFeatureFlagInterface
	clockService       clock.IClock
	cache              *cache
}

func NewFeatureFlagWithCacheService(errorLoggerService errorlogger.ErrorLogger, clockService clock.IClock, ttl int64) ServiceFeatureFlagInterface {
	featureFlagService := NewFeatureFlagService(errorLoggerService)

	cacheHandler := &cache{
		ttl:          ttl,
		clockService: clockService,
	}
	cachedService := &serviceFeatureFlagWithCache{
		featureFlagService: featureFlagService,
		clockService:       clockService,
		cache:              cacheHandler,
	}

	return cachedService
}

func (s *serviceFeatureFlagWithCache) IsActive(ormService *beeorm.Engine, name string) bool {
	featureFlagEntity, has := s.getFeatureFlagEntity(ormService, name)
	if !has {
		return false
	}

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlagWithCache) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return fmt.Errorf("feature (%s) is not active", name)
	}

	return nil
}

func (s *serviceFeatureFlagWithCache) Enable(ormService *beeorm.Engine, name string) error {
	err := s.featureFlagService.Enable(ormService, name)
	s.cache.Delete(name)

	return err
}

func (s *serviceFeatureFlagWithCache) Disable(ormService *beeorm.Engine, name string) error {
	err := s.featureFlagService.Disable(ormService, name)
	s.cache.Delete(name)

	return err
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

func (s *serviceFeatureFlagWithCache) getFeatureFlagEntity(ormService *beeorm.Engine, name string) (*entity.FeatureFlagEntity, bool) {
	if name == "" {
		panic("name cannot be empty")
	}

	if featureFlagEntity, has := s.cache.Load(name); has {
		return featureFlagEntity, true
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}

	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return nil, false
	}

	s.cache.Store(name, cacheEntry{
		featureFlagEntity: featureFlagEntity,
		time:              s.clockService.Now(),
	})

	return featureFlagEntity, true
}
