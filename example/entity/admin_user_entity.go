package entity

import (
	"github.com/coretrix/trixorm"

	hitrixEntity "github.com/coretrix/hitrix/pkg/entity"
)

type AdminUserEntity struct {
	trixorm.ORM `orm:"table=admin_users;log=log_db_pool;redisCache;redisSearch=search_pool"`
	ID          uint64
	RoleID      *hitrixEntity.RoleEntity `orm:"required"`
}

func (u *AdminUserEntity) SetRole(roleEntity *hitrixEntity.RoleEntity) {
	u.RoleID = roleEntity
}
