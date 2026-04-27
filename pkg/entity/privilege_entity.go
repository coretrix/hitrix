package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type PrivilegeEntity struct {
	trixorm.ORM   `orm:"table=privileges;redisCache"`
	ID            uint64
	RoleID        *RoleEntity         `orm:"required;unique=RoleID_ResourceID_FakeDelete:1"`
	ResourceID    *ResourceEntity     `orm:"required;unique=RoleID_ResourceID_FakeDelete:2"`
	PermissionIDs []*PermissionEntity `orm:"required"`
	CreatedAt     time.Time           `orm:"time=true"`
	FakeDelete    bool                `orm:"unique=RoleID_ResourceID_FakeDelete:3"`

	CachedQueryPrivilegeRoleIDResourceID *trixorm.CachedQuery `queryOne:":RoleID = ? AND :ResourceID = ?"`
	CachedQueryPrivilegeRoleID           *trixorm.CachedQuery `query:":RoleID = ?"`
	CachedQueryPrivilegeResourceID       *trixorm.CachedQuery `query:":ResourceID = ?"`
}
