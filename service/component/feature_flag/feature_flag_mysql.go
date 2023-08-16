package featureflag

import (
	"errors"
	"fmt"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
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

func (s *serviceFeatureFlag) IsActive(ormService *datalayer.DataLayer, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}

	found := ormService.RedisSearchOne(featureFlagEntity, query)
	if !found {
		return false
	}

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlag) FailIfIsNotActive(ormService *datalayer.DataLayer, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return fmt.Errorf("feature (%s) is not active", name)
	}

	return nil
}

func (s *serviceFeatureFlag) Enable(ormService *datalayer.DataLayer, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := redisearch.NewRedisSearchQuery()
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

func (s *serviceFeatureFlag) Disable(ormService *datalayer.DataLayer, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	query := redisearch.NewRedisSearchQuery()
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

func (s *serviceFeatureFlag) getAllActive(ormService *datalayer.DataLayer, pager *beeorm.Pager) []IFeatureFlag {
	query := redisearch.NewRedisSearchQuery()
	query.FilterBool("Registered", true)
	query.FilterBool("Enabled", true)

	var featureFlagEntities []*entity.FeatureFlagEntity
	ormService.RedisSearchMany(query, pager, &featureFlagEntities)

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

func (s *serviceFeatureFlag) GetScriptsSingleInstance(ormService *datalayer.DataLayer) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, beeorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsSingleInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) GetScriptsMultiInstance(ormService *datalayer.DataLayer) []app.IScript {
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
		if _, has := s.featureFlags[featureFlag.GetName()]; has {
			panic("feature flag with name '" + featureFlag.GetName() + "' already exists")
		}

		s.featureFlags[featureFlag.GetName()] = featureFlag
	}
}

func (s *serviceFeatureFlag) Sync(ormService *datalayer.DataLayer, clockService clock.IClock) {
	var featureFlagEntities []*entity.FeatureFlagEntity
	var lastID uint64

	for {
		var rows []*entity.FeatureFlagEntity
		pager := beeorm.NewPager(1, 1000)
		ormService.Search(beeorm.NewWhere("ID > ? ORDER BY ID ASC", lastID), pager, &rows)

		if len(rows) == 0 {
			break
		}

		lastID = rows[len(rows)-1].ID
		featureFlagEntities = append(featureFlagEntities, rows...)

		if len(rows) < pager.PageSize {
			break
		}
	}

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

			err := ormService.FlushWithCheck(featureFlagEntity)

			if err != nil {
				if duplicateKeyError, ok := err.(*beeorm.DuplicatedKeyError); ok {
					if duplicateKeyError.Index != "Name" {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
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
