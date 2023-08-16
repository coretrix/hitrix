package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type RoleEntity struct {
	beeorm.ORM   `orm:"table=roles;redisCache;redisSearch=search_pool"`
	ID           uint64    `orm:"sortable"`
	Name         string    `orm:"required;searchable;unique=Name_FakeDelete:1"`
	IsPredefined bool      `orm:"searchable"`
	CreatedAt    time.Time `orm:"time=true"`
	FakeDelete   bool      `orm:"unique=Name_FakeDelete:2"`
}
