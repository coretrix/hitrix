package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type PermissionEntity struct {
	trixorm.ORM `orm:"table=permissions;redisCache"`
	ID          uint64
	ResourceID  *ResourceEntity `orm:"required;unique=ResourceID_Name_FakeDelete:1"`
	Name        string          `orm:"required;unique=ResourceID_Name_FakeDelete:3"`
	CreatedAt   time.Time       `orm:"time=true"`
	FakeDelete  bool            `orm:"unique=ResourceID_Name_FakeDelete:2"`

	CachedQueryAll        *trixorm.CachedQuery `query:"1 ORDER BY ID"`
	CachedQueryResourceID *trixorm.CachedQuery `query:":ResourceID = ?"`
}
