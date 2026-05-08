package entity

import (
	"github.com/coretrix/trixorm"
)

type DevPanelUserEntity struct {
	trixorm.ORM `orm:"table=dev_panel_users;redisCache;redisSearch=search_pool"`
	ID          uint64
	Email       string `orm:"unique=Email;searchable=text"`
	Password    string
}

func (u *DevPanelUserEntity) GetUniqueFieldName() string {
	return "Email"
}

func (u *DevPanelUserEntity) GetUsername() string {
	return u.Email
}

func (u *DevPanelUserEntity) GetPassword() string {
	return u.Password
}

func (u *DevPanelUserEntity) CanAuthenticate() bool {
	return true
}
