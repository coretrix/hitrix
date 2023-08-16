package entity

import "github.com/latolukasz/beeorm/v2"

type DevPanelUserEntity struct {
	beeorm.ORM `orm:"crud-stream;table=dev_panel_users;redisCache;redisSearch=search_pool"`
	ID         uint64
	Email      string `orm:"unique=Email;searchable"`
	Password   string `orm:"searchable"`
	FakeDelete bool
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
