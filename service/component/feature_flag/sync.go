package featureflag

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

func syncFeatureFlags(
	ormService *beeorm.Engine,
	clockService clock.IClock,
	errorLogger errorlogger.ErrorLogger,
	featureFlags map[string]IFeatureFlag,
) {
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
			errorLogger.LogError("feature flag is nil")
		}
	}

	for _, registeredFeatureFlag := range featureFlags {
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
		if _, ok := featureFlags[name]; !ok {
			dbFeatureFlag.Registered = false
			dbFeatureFlag.UpdatedAt = clockService.NowPointer()

			flusher.Track(dbFeatureFlag)
		}
	}

	flusher.Flush()
}
