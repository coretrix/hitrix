package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type PrivilegeEntity struct {
	beeorm.ORM    `orm:"table=privileges;redisCache"`
	ID            uint64
	RoleID        *RoleEntity         `orm:"required;unique=RoleID_ResourceID_FakeDelete:1"`
	ResourceID    *ResourceEntity     `orm:"required;unique=RoleID_ResourceID_FakeDelete:2"`
	PermissionIDs []*PermissionEntity `orm:"required"`
	CreatedAt     time.Time           `orm:"time=true"`
	FakeDelete    bool                `orm:"unique=RoleID_ResourceID_FakeDelete:3"`

	CachedQueryPrivilegeRoleIDResourceID *beeorm.CachedQuery `queryOne:":RoleID = ? AND :ResourceID = ?"`
	CachedQueryPrivilegeRoleID           *beeorm.CachedQuery `query:":RoleID = ?"`
	CachedQueryPrivilegeResourceID       *beeorm.CachedQuery `query:":ResourceID = ?"`
}
