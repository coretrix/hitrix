package entity

import (
	"github.com/latolukasz/beeorm"
)

type DevPanelUserEntity struct {
	beeorm.ORM `orm:"table=dev_panel_users;redisCache;redisSearch=search"`
	ID         uint64
	Email      string `orm:"unique=Email;searchable"`
	Password   string
}

func (e *DevPanelUserEntity) GetUniqueFieldName() string {
	return "Email"
}

func (e *DevPanelUserEntity) GetUsername() string {
	return e.Email
}

func (e *DevPanelUserEntity) GetPassword() string {
	return e.Password
}
