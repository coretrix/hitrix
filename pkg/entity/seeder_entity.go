package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type SeederEntity struct {
	beeorm.ORM `orm:"table=seeder;redisCache;redisSearch=search_pool;"`
	ID         uint64
	Name       string    `orm:"required;unique=Seeder_Name;searchable;"`
	CreatedAt  time.Time `orm:"time=true"`
}
