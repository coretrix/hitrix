package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type ResourceEntity struct {
	trixorm.ORM `orm:"table=resources;redisCache;redisSearch=search_pool"`
	ID          uint64    `orm:"searchable"`
	Name        string    `orm:"required;searchable=text;unique=Name_FakeDelete:1"`
	CreatedAt   time.Time `orm:"time=true"`
	FakeDelete  bool      `orm:"unique=Name_FakeDelete:2"`

	CachedQueryName *trixorm.CachedQuery `queryOne:":Name = ?"`
}
