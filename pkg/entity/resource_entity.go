package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type ResourceEntity struct {
	beeorm.ORM `orm:"table=resources;redisCache;redisSearch=search_pool"`
	ID         uint64    `orm:"searchable"`
	Name       string    `orm:"required;searchable;unique=Name_FakeDelete:1"`
	CreatedAt  time.Time `orm:"time=true"`
	FakeDelete bool      `orm:"unique=Name_FakeDelete:2"`
}
