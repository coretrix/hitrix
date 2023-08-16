package setting

import (
	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

type ServiceSettingInterface interface {
	Get(ormService *datalayer.DataLayer, key string) (*entity.SettingsEntity, bool)
	GetString(ormService *datalayer.DataLayer, key string) (string, bool)
	GetInt(ormService *datalayer.DataLayer, key string) (int, bool)
	GetUint(ormService *datalayer.DataLayer, key string) (uint, bool)
	GetInt64(ormService *datalayer.DataLayer, key string) (int64, bool)
	GetUint64(ormService *datalayer.DataLayer, key string) (uint64, bool)
	GetFloat64(ormService *datalayer.DataLayer, key string) (float64, bool)
	GetBool(ormService *datalayer.DataLayer, key string) (bool, bool)
}
