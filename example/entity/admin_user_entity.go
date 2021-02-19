package entity

import (
	"github.com/summer-solutions/orm"
)

type AdminUserEntity struct {
	orm.ORM  `orm:"table=admin_users;redisCache"`
	ID       uint64
	Email    string `orm:"unique=Email"`
	Password string

	UserEmailIndex *orm.CachedQuery `queryOne:":Email = ?"`
}

func (e *AdminUserEntity) GetUsername() string {
	return e.Email
}

func (e *AdminUserEntity) GetPassword() string {
	return e.Password
}
