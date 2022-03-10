package entity

import (
	"github.com/latolukasz/beeorm"
)

const (
	SettingsValueTypeText     = "text"
	SettingsValueTypeNumber   = "number"
	SettingsValueTypeEmail    = "email"
	SettingsValueTypeTel      = "tel"
	SettingsValueTypeURL      = "url"
	SettingsValueTypePassword = "password"
)

type settingsValueType struct {
	SettingsValueTypeText     string
	SettingsValueTypeNumber   string
	SettingsValueTypeEmail    string
	SettingsValueTypeTel      string
	SettingsValueTypeURL      string
	SettingsValueTypePassword string
}

var SettingsValueTypeAll = settingsValueType{
	SettingsValueTypeText:     SettingsValueTypeText,
	SettingsValueTypeNumber:   SettingsValueTypeNumber,
	SettingsValueTypeEmail:    SettingsValueTypeEmail,
	SettingsValueTypeTel:      SettingsValueTypeTel,
	SettingsValueTypeURL:      SettingsValueTypeURL,
	SettingsValueTypePassword: SettingsValueTypePassword,
}

type SettingsEntity struct {
	beeorm.ORM `orm:"table=settings;redisCache;redisSearch=search_pool;"`
	ID         uint64 `orm:"sortable"`
	Key        string `orm:"required;unique=Settings_Key;sortable;searchable;"`
	Value      string `orm:"required;length=max;"`
	ValueType  string `orm:"enum=entity.SettingsValueTypeAll"`
	Editable   bool
	Deletable  bool
	Hidden     bool `orm:"searchable"`
}
