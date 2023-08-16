package setting

import (
	"strconv"
	"strings"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

type serviceSetting struct {
	cache map[string]*entity.SettingsEntity
}

func NewSettingService() ServiceSettingInterface {
	return &serviceSetting{cache: map[string]*entity.SettingsEntity{}}
}

func (s *serviceSetting) Get(ormService *datalayer.ORM, key string) (*entity.SettingsEntity, bool) {
	if cachedEntity, exists := s.cache[key]; exists {
		return cachedEntity, true
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterString("Key", key)

	settingEntity := &entity.SettingsEntity{}

	found := ormService.RedisSearchOne(settingEntity, query)
	if !found {
		return nil, false
	}

	if !settingEntity.Editable && !settingEntity.Deletable {
		s.cache[key] = settingEntity
	}

	return settingEntity, true
}

func (s *serviceSetting) GetString(ormService *datalayer.ORM, key string) (string, bool) {
	setting, found := s.Get(ormService, key)
	if found {
		return setting.Value, true
	}

	return "", false
}

func (s *serviceSetting) GetInt(ormService *datalayer.ORM, key string) (int, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return int(i), true
}

func (s *serviceSetting) GetUint(ormService *datalayer.ORM, key string) (uint, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseUint(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(i), true
}

func (s *serviceSetting) GetInt64(ormService *datalayer.ORM, key string) (int64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetUint64(ormService *datalayer.ORM, key string) (uint64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseUint(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetFloat64(ormService *datalayer.ORM, key string) (float64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetBool(ormService *datalayer.ORM, key string) (bool, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return false, false
	}

	if strings.ToLower(setting.Value) == "false" {
		return false, true
	}

	return true, true
}
