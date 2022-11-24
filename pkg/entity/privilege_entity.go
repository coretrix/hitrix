package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type PrivilegeEntity struct {
	beeorm.ORM    `orm:"table=privileges;redisCache;redisSearch=search_pool"`
	ID            uint64
	RoleID        *RoleEntity         `orm:"required;searchable;unique=RoleID_ResourceID_FakeDelete:1"`
	ResourceID    *ResourceEntity     `orm:"required;searchable;unique=RoleID_ResourceID_FakeDelete:2"`
	PermissionIDs []*PermissionEntity `orm:"required;searchable"`
	CreatedAt     time.Time           `orm:"time=true"`
	FakeDelete    bool                `orm:"unique=RoleID_ResourceID_FakeDelete:3"`
}
