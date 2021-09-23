package entity

import (
	"github.com/latolukasz/beeorm"
)

const (
	seedsSetting      = "seeds"
	trustpilotSetting = "trustpilot-accesstoken"
)

type HitrixSettings struct {
	Seeds      string
	Trustpilot string
}

var HitrixSettingAll = HitrixSettings{
	Seeds:      seedsSetting,
	Trustpilot: trustpilotSetting,
}

type SettingSeedsValue map[string]int

type SettingsEntity struct {
	beeorm.ORM `orm:"table=settings;redisCache;redisSearch=search_pool;"`
	ID         uint64
	Key        string `orm:"required;unique=Settings_Key;searchable;"`
	Value      string `orm:"required;length=max;"`
}
