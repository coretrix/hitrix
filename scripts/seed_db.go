package scripts

import (
	"context"
	"encoding/json"

	"github.com/coretrix/hitrix"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

type DBSeedScript struct {
	Seeds map[string]Seed
}

func (script *DBSeedScript) Run(ctx context.Context, _ hitrix.Exit) {
	ormService, _ := service.DI().OrmEngine()
	Seeder(script.Seeds, ormService)
}

func (script *DBSeedScript) Unique() bool {
	return true
}

func (script *DBSeedScript) Description() string {
	return "Seed Database"
}

type Seed interface {
	Execute(*beeorm.Engine)
	Version() int
}

func Seeder(seeds map[string]Seed, ormService *beeorm.Engine) {

	var setting entity.SettingsEntity

	whereStmt := beeorm.NewWhere("`Key` = ?", entity.HitrixSettingAll.Seeds)
	var hasExecutedSeedsSetting = ormService.SearchOne(whereStmt, &setting)

	var executedSeeds entity.SettingSeedsValue
	if hasExecutedSeedsSetting {
		if err := json.Unmarshal([]byte(setting.Value), &executedSeeds); err != nil {
			panic(err.Error())
		}
	}

	var newSeeds entity.SettingSeedsValue = make(entity.SettingSeedsValue)

	for k, seed := range seeds {
		_, hasExecutedSeed := executedSeeds[k]
		if !hasExecutedSeedsSetting || !hasExecutedSeed ||
			(hasExecutedSeed && executedSeeds[k] < seed.Version()) {
			seed.Execute(ormService)
			newSeeds[k] = seed.Version()
		}
	}
	if len(newSeeds) > 0 {
		saveNewSeeds(ormService, newSeeds)
	}

}
func saveNewSeeds(ormService *beeorm.Engine, newSeeds entity.SettingSeedsValue) {
	settingsEntity := &entity.SettingsEntity{}

	query := &beeorm.RedisSearchQuery{}
	query.FilterString("Key", entity.HitrixSettingAll.Seeds)

	hasExecutedSeedsSetting:= ormService.RedisSearchOne(settingsEntity, query)

	if hasExecutedSeedsSetting {
		var oldSeeds entity.SettingSeedsValue
		if err := json.Unmarshal([]byte(settingsEntity.Value), &oldSeeds); err != nil {
			panic(err.Error())
		}

		// overwrite old with newSeeds
		for k, v := range newSeeds {
			oldSeeds[k] = v
		}
		newSeeds = oldSeeds
	}else {
		settingsEntity.Key = entity.HitrixSettingAll.Seeds
	}

	str, err := json.Marshal(newSeeds)
	if err != nil {
		panic(err.Error())
	}
	settingsEntity.Value = string(str)
	ormService.Flush(settingsEntity)
}
