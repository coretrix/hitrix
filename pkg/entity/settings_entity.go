package entity

import "github.com/latolukasz/beeorm"

const (
	seedsSetting = "seeds"
)

type HitrixSettings struct {
	Seeds string
}

var HitrixSettingAll = HitrixSettings{
	Seeds: seedsSetting,
}

type SettingSeedsValue map[string]int

type SettingsEntity struct {
	beeorm.ORM `orm:"table=settings;redisCache;redisSearch=search_pool;"`
	ID         uint64
	Key        string `orm:"required;unique=Settings_Key;searchable;"`
	Value      string `orm:"required;length=max;"`
}
