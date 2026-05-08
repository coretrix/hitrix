package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type RoleEntity struct {
	trixorm.ORM  `orm:"table=roles;redisCache;redisSearch=search_pool"`
	ID           uint64    `orm:"sortable"`
	Name         string    `orm:"required;searchable=text;unique=Name_FakeDelete:1"`
	IsPredefined bool      `orm:"searchable"`
	CreatedAt    time.Time `orm:"time=true"`
	FakeDelete   bool      `orm:"unique=Name_FakeDelete:2"`
}
