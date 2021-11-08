package featureflag

import (
	"errors"

	"github.com/coretrix/hitrix/service/component/app"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

type IFeatureFlag interface {
	GetName() string
	ScriptsSingleInstance() []app.IScript
	ScriptsMultiInstance() []app.IScript
}

type serviceFeatureFlag struct {
	featureFlags       map[string]IFeatureFlag
	errorLoggerService errorlogger.ErrorLogger
}

func NewFeatureFlagService(errorLoggerService errorlogger.ErrorLogger) ServiceFeatureFlagInterface {
	return &serviceFeatureFlag{
		errorLoggerService: errorLoggerService,
	}
}

func (s *serviceFeatureFlag) IsActive(ormService *beeorm.Engine, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return false
	}

	return featureFlagEntity.IsActive
}

func (s *serviceFeatureFlag) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return errors.New("feature is not active")
	}

	return nil
}

func (s *serviceFeatureFlag) Activate(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	featureFlagEntity.IsActive = true
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) DeActivate(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	featureFlagEntity.IsActive = false
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) Create(ormService *beeorm.Engine, clockService clock.IClock, name string, isActive bool) error {
	if name == "" {
		panic("name cannot be empty")
	}

	featureFlagEntity := &entity.FeatureFlagEntity{
		Name:      name,
		IsActive:  isActive,
		UpdatedAt: nil,
		CreatedAt: clockService.Now(),
	}

	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) Delete(ormService *beeorm.Engine, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Name", name)
	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return errors.New("feature cannot be found")
	}

	ormService.Delete(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) GetAll(ormService *beeorm.Engine, pager *beeorm.Pager) []*entity.FeatureFlagEntity {
	query := beeorm.NewRedisSearchQuery()
	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, pager)

	return featureFlagEntities
}

func (s *serviceFeatureFlag) getAllActive(ormService *beeorm.Engine, pager *beeorm.Pager) []IFeatureFlag {
	query := beeorm.NewRedisSearchQuery()
	query.FilterBool("IsActive", true)

	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, pager)

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

func (s *serviceFeatureFlag) GetScriptsSingleInstance(ormService *beeorm.Engine) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, beeorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsSingleInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) GetScriptsMultiInstance(ormService *beeorm.Engine) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, beeorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsSingleInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) Register(featureFlags ...IFeatureFlag) {
	s.featureFlags = map[string]IFeatureFlag{}

	for _, featureFlag := range featureFlags {
		s.featureFlags[featureFlag.GetName()] = featureFlag
	}
}

func (s *serviceFeatureFlag) Sync(ormService *beeorm.Engine, clockService clock.IClock) {
	query := beeorm.NewRedisSearchQuery()

	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, beeorm.NewPager(1, 1000))

	dbFeatureFlags := map[string]struct{}{}
	for _, featureFlagEntity := range featureFlagEntities {
		dbFeatureFlags[featureFlagEntity.Name] = struct{}{}
	}

	for _, registeredFeatureFlag := range s.featureFlags {
		if _, ok := dbFeatureFlags[registeredFeatureFlag.GetName()]; !ok {
			err := s.Create(ormService, clockService, registeredFeatureFlag.GetName(), false)
			if err != nil {
				s.errorLoggerService.LogError(err)
			}
		}
	}
}
