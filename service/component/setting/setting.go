package setting

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type ServiceSettingInterface interface {
	Get(ormService *beeorm.Engine, key string) (*entity.SettingsEntity, bool)
	GetString(ormService *beeorm.Engine, key string) (string, bool)
	GetInt(ormService *beeorm.Engine, key string) (int, bool)
	GetUint(ormService *beeorm.Engine, key string) (uint, bool)
	GetInt64(ormService *beeorm.Engine, key string) (int64, bool)
	GetUint64(ormService *beeorm.Engine, key string) (uint64, bool)
	GetFloat64(ormService *beeorm.Engine, key string) (float64, bool)
	GetBool(ormService *beeorm.Engine, key string) (bool, bool)
}
