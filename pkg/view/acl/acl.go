package acl

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

func ACL(ormService *beeorm.Engine, roleEntity *entity.RoleEntity, resource string, permissions ...string) bool {
	resourceQuery := beeorm.NewRedisSearchQuery()
	resourceQuery.FilterString("Name", resource)

	resourceEntity := &entity.ResourceEntity{}
	if !ormService.RedisSearchOne(resourceEntity, resourceQuery) {
		return false
	}

	permissionQuery := beeorm.NewRedisSearchQuery()
	permissionQuery.FilterUint("ResourceID", resourceEntity.ID)
	permissionQuery.FilterString("Name", permissions...)

	permissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(&permissionEntities, permissionQuery, beeorm.NewPager(1, 1000))

	if len(permissions) != len(permissionEntities) {
		return false
	}

	permissionIDs := make([]uint64, len(permissionEntities))

	for i, permissionEntity := range permissionEntities {
		permissionIDs[i] = permissionEntity.ID
	}

	privilegeQuery := beeorm.NewRedisSearchQuery()
	privilegeQuery.FilterUint("RoleID", roleEntity.ID)
	privilegeQuery.FilterUint("ResourceID", resourceEntity.ID)
	privilegeQuery.FilterManyReferenceIn("PermissionIDs", permissionIDs...)

	return ormService.RedisSearchOne(&entity.PrivilegeEntity{}, privilegeQuery)
}
