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
	featureFlags := make(map[string]IFeatureFlag)
	return &serviceFeatureFlag{
		featureFlags:       featureFlags,
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

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlag) FailIfIsNotActive(ormService *beeorm.Engine, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return errors.New("feature is not active")
	}

	return nil
}

func (s *serviceFeatureFlag) Enable(ormService *beeorm.Engine, name string) error {
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

	featureFlagEntity.Enabled = true
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) Disable(ormService *beeorm.Engine, name string) error {
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

	featureFlagEntity.Enabled = false
	ormService.Flush(featureFlagEntity)

	return nil
}

func (s *serviceFeatureFlag) getAllActive(ormService *beeorm.Engine, pager *beeorm.Pager) []IFeatureFlag {
	query := beeorm.NewRedisSearchQuery()
	query.FilterBool("Registered", true)
	query.FilterBool("Enabled", true)

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
		allScripts = append(allScripts, featureFlag.ScriptsMultiInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) Register(featureFlags ...IFeatureFlag) {
	s.featureFlags = make(map[string]IFeatureFlag)

	for _, featureFlag := range featureFlags {
		s.featureFlags[featureFlag.GetName()] = featureFlag
	}
}

func (s *serviceFeatureFlag) Sync(ormService *beeorm.Engine, clockService clock.IClock) {
	query := beeorm.NewRedisSearchQuery()

	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearch(&featureFlagEntities, query, beeorm.NewPager(1, 1000))
	flusher := ormService.NewFlusher()

	dbFeatureFlags := make(map[string]*entity.FeatureFlagEntity)

	for _, featureFlagEntity := range featureFlagEntities {
		if featureFlagEntity != nil {
			dbFeatureFlags[featureFlagEntity.Name] = featureFlagEntity
		} else {
			s.errorLoggerService.LogError("feature flag is nil")
		}
	}

	for _, registeredFeatureFlag := range s.featureFlags {
		if _, ok := dbFeatureFlags[registeredFeatureFlag.GetName()]; !ok {
			featureFlagEntity := &entity.FeatureFlagEntity{
				Name:       registeredFeatureFlag.GetName(),
				Registered: true,
				Enabled:    false,
				UpdatedAt:  nil,
				CreatedAt:  clockService.Now(),
			}

			flusher.Track(featureFlagEntity)
		} else if !dbFeatureFlags[registeredFeatureFlag.GetName()].Registered {
			dbFeatureFlags[registeredFeatureFlag.GetName()].Registered = true
			dbFeatureFlags[registeredFeatureFlag.GetName()].UpdatedAt = clockService.NowPointer()

			flusher.Track(dbFeatureFlags[registeredFeatureFlag.GetName()])
		}
	}

	for name, dbFeatureFlag := range dbFeatureFlags {
		if _, ok := s.featureFlags[name]; !ok {
			dbFeatureFlag.Registered = false
			dbFeatureFlag.UpdatedAt = clockService.NowPointer()

			flusher.Track(dbFeatureFlag)
		}
	}

	flusher.Flush()
}