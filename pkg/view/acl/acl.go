package acl

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

func ACL(ormService *datalayer.ORM, roleEntity *entity.RoleEntity, resource string, permissions ...string) bool {
	resourceQuery := redisearch.NewRedisSearchQuery()
	resourceQuery.FilterString("Name", resource)

	resourceEntity := &entity.ResourceEntity{}
	if !ormService.RedisSearchOne(resourceEntity, resourceQuery) {
		return false
	}

	permissionQuery := redisearch.NewRedisSearchQuery()
	permissionQuery.FilterUint("ResourceID", resourceEntity.ID)
	permissionQuery.FilterString("Name", permissions...)

	permissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(permissionQuery, beeorm.NewPager(1, 1000), &permissionEntities)

	if len(permissions) != len(permissionEntities) {
		return false
	}

	permissionIDs := make([]uint64, len(permissionEntities))

	for i, permissionEntity := range permissionEntities {
		permissionIDs[i] = permissionEntity.ID
	}

	privilegeQuery := redisearch.NewRedisSearchQuery()
	privilegeQuery.FilterUint("RoleID", roleEntity.ID)
	privilegeQuery.FilterUint("ResourceID", resourceEntity.ID)
	privilegeQuery.FilterManyReferenceIn("PermissionIDs", permissionIDs...)

	return ormService.RedisSearchOne(&entity.PrivilegeEntity{}, privilegeQuery)
}
