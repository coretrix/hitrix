package entity

import (
	"github.com/latolukasz/orm"
)

type AdminUserEntity struct {
	orm.ORM  `orm:"table=admin_users;redisCache;redisSearch=search"`
	ID       uint64
	Email    string `orm:"unique=Email;searchable"`
	Password string
}

func (e *AdminUserEntity) GetUniqueFieldName() string {
	return "Email"
}

func (e *AdminUserEntity) GetUsername() string {
	return e.Email
}

func (e *AdminUserEntity) GetPassword() string {
	return e.Password
}
