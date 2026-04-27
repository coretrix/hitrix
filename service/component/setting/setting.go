package setting

import (
	"github.com/coretrix/trixorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type ServiceSettingInterface interface {
	Get(ormService *trixorm.Engine, key string) (*entity.SettingsEntity, bool)
	GetString(ormService *trixorm.Engine, key string) (string, bool)
	GetInt(ormService *trixorm.Engine, key string) (int, bool)
	GetUint(ormService *trixorm.Engine, key string) (uint, bool)
	GetInt64(ormService *trixorm.Engine, key string) (int64, bool)
	GetUint64(ormService *trixorm.Engine, key string) (uint64, bool)
	GetFloat64(ormService *trixorm.Engine, key string) (float64, bool)
	GetBool(ormService *trixorm.Engine, key string) (bool, bool)
}
