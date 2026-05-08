package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type SeederEntity struct {
	trixorm.ORM `orm:"table=seeder;redisCache;redisSearch=search_pool;"`
	ID          uint64
	Name        string    `orm:"required;unique=Seeder_Name;searchable=text"`
	CreatedAt   time.Time `orm:"time=true"`
}
