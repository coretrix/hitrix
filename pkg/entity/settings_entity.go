package entity

import (
	"github.com/latolukasz/beeorm"
)

type SettingsEntity struct {
	beeorm.ORM `orm:"table=settings;redisCache;redisSearch=search_pool;"`
	ID         uint64 `orm:"sortable"`
	Key        string `orm:"required;unique=Settings_Key;sortable;searchable;"`
	Value      string `orm:"required;length=max;"`
	Editable   bool
	Deletable  bool
	Hidden     bool `orm:"searchable"`
}
