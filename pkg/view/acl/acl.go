package acl

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

func ACL(ormService *beeorm.Engine, roleEntity *entity.RoleEntity, resource string, permissions ...string) bool {
	resourceEntity := &entity.ResourceEntity{}
	if !ormService.CachedSearchOne(resourceEntity, "CachedQueryName", resource) {
		return false
	}

	allPermissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.CachedSearch(
		&allPermissionEntities,
		"CachedQueryResourceID",
		beeorm.NewPager(1, 1000),
		resourceEntity.ID,
	)

	permissionEntities := make([]*entity.PermissionEntity, 0)

	for _, permissionEntity := range allPermissionEntities {
		for _, permission := range permissions {
			if permissionEntity.Name == permission {
				permissionEntities = append(permissionEntities, permissionEntity)
			}
		}
	}

	if len(permissions) != len(permissionEntities) {
		return false
	}

	permissionIDs := make([]uint64, len(permissionEntities))

	for i, permissionEntity := range permissionEntities {
		permissionIDs[i] = permissionEntity.ID
	}

	privilegeEntity := &entity.PrivilegeEntity{}
	ormService.CachedSearchOne(
		privilegeEntity,
		"CachedQueryPrivilegeRoleIDResourceID",
		roleEntity.ID,
		resourceEntity.ID,
	)

	hasPrivilege := false

	for _, permissionEntity := range privilegeEntity.PermissionIDs {
		for _, permissionID := range permissionIDs {
			if permissionEntity.ID == permissionID {
				hasPrivilege = true

				break
			}
		}
	}

	return hasPrivilege
}
