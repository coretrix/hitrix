package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type PermissionEntity struct {
	beeorm.ORM `orm:"table=permissions;redisCache;redisSearch=search_pool"`
	ID         uint64          `orm:"searchable;sortable"`
	ResourceID *ResourceEntity `orm:"required;searchable;unique=ResourceID_Name_FakeDelete:1"`
	Name       string          `orm:"required;searchable;unique=ResourceID_Name_FakeDelete:3"`
	CreatedAt  time.Time       `orm:"time=true"`
	FakeDelete bool            `orm:"unique=ResourceID_Name_FakeDelete:2"`

	CachedQueryAll        *beeorm.CachedQuery `query:"1 ORDER BY ID"`
	CachedQueryResourceID *beeorm.CachedQuery `query:":ResourceID = ?"`
}
