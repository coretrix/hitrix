package featureflag

import (
	"errors"
	"fmt"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type serviceFeatureFlagMysql struct {
	featureFlags       map[string]IFeatureFlag
	errorLoggerService errorlogger.ErrorLogger
	clockService       clock.IClock
	cache              *cache
}

func NewFeatureFlagMysqlOnlyService(errorLoggerService errorlogger.ErrorLogger, clockService clock.IClock, ttl int64) ServiceFeatureFlagInterface {
	featureFlags := make(map[string]IFeatureFlag)

	cacheHandler := &cache{
		ttl:          ttl,
		clockService: clockService,
	}

	return &serviceFeatureFlagMysql{
		featureFlags:       featureFlags,
		errorLoggerService: errorLoggerService,
		clockService:       clockService,
		cache:              cacheHandler,
	}
}

func (s *serviceFeatureFlagMysql) IsActive(ormService *beeorm.Engine, name string) bool {
	featureFlagEntity, found := s.getFeatureFlagEntity(ormService, name)
	if !found {
		return false
	}

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlagMysql) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	if isActive := s.IsActive(ormService, name); !isActive {
		return fmt.Errorf("feature (%s) is not active", name)
	}

	return nil
}

func (s *serviceFeatureFlagMysql) Enable(ormService *beeorm.Engine, name string) error {
	featureFlagEntity, found := s.getFeatureFlagEntity(ormService, name)

	if !found {
		return errors.New("feature not found")
	}

	featureFlagEntity.Enabled = true
	ormService.Flush(featureFlagEntity)
	s.cache.Delete(name)

	return nil
}

func (s *serviceFeatureFlagMysql) Disable(ormService *beeorm.Engine, name string) error {
	featureFlagEntity, found := s.getFeatureFlagEntity(ormService, name)

	if !found {
		return errors.New("feature cannot be found")
	}

	featureFlagEntity.Enabled = false
	ormService.Flush(featureFlagEntity)
	s.cache.Delete(name)

	return nil
}

func (s *serviceFeatureFlagMysql) GetScriptsSingleInstance(ormService *beeorm.Engine) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, beeorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsSingleInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlagMysql) GetScriptsMultiInstance(ormService *beeorm.Engine) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, beeorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsMultiInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlagMysql) Register(featureFlags ...IFeatureFlag) {
	s.featureFlags = make(map[string]IFeatureFlag)

	for _, featureFlag := range featureFlags {
		if _, has := s.featureFlags[featureFlag.GetName()]; has {
			panic("feature flag with name '" + featureFlag.GetName() + "' already exists")
		}

		s.featureFlags[featureFlag.GetName()] = featureFlag
	}
}
func (s *serviceFeatureFlagMysql) Sync(ormService *beeorm.Engine, clockService clock.IClock) {
	syncFeatureFlags(ormService, clockService, s.errorLoggerService, s.featureFlags)
}

func (s *serviceFeatureFlagMysql) getFeatureFlagEntity(ormService *beeorm.Engine, name string) (*entity.FeatureFlagEntity, bool) {
	if name == "" {
		panic("name cannot be empty")
	}

	if featureFlagEntity, has := s.cache.Load(name); has {
		return featureFlagEntity, true
	}

	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.SearchOne(beeorm.NewWhere("Name = ?", name), featureFlagEntity)

	if !found {
		return nil, false
	}

	s.cache.Store(name, cacheEntry{
		featureFlagEntity: featureFlagEntity,
		time:              s.clockService.Now(),
	})

	return featureFlagEntity, true
}
func (s *serviceFeatureFlagMysql) getAllActive(ormService *beeorm.Engine, pager *beeorm.Pager) []IFeatureFlag {
	where := beeorm.NewWhere("Registered = ? AND Enabled = ?", true, true)

	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.Search(where, pager, &featureFlagEntities)

	activeFeatureFlags := make([]IFeatureFlag, 0)

	for _, featureFlagEntity := range featureFlagEntities {
		if _, ok := s.featureFlags[featureFlagEntity.Name]; !ok {
			s.errorLoggerService.LogError("feature flag " + featureFlagEntity.Name + " is not registered")

			continue
		}

		activeFeatureFlags = append(activeFeatureFlags, s.featureFlags[featureFlagEntity.Name])
	}

	return activeFeatureFlags
}
