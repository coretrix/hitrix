package setting

import (
	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

type ServiceSettingInterface interface {
	Get(ormService *datalayer.ORM, key string) (*entity.SettingsEntity, bool)
	GetString(ormService *datalayer.ORM, key string) (string, bool)
	GetInt(ormService *datalayer.ORM, key string) (int, bool)
	GetUint(ormService *datalayer.ORM, key string) (uint, bool)
	GetInt64(ormService *datalayer.ORM, key string) (int64, bool)
	GetUint64(ormService *datalayer.ORM, key string) (uint64, bool)
	GetFloat64(ormService *datalayer.ORM, key string) (float64, bool)
	GetBool(ormService *datalayer.ORM, key string) (bool, bool)
}
