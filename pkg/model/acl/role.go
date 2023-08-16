package acl

import (
	"fmt"
	"time"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service"
)

type UserRoleSetter interface {
	SetRole(roleEntity *entity.RoleEntity)
}

func CreateRole(c *gin.Context, request *acl.CreateOrUpdateRoleRequestDTO) error {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
	if err != nil {
		return err
	}

	now := service.DI().Clock().Now()

	err = helper.DBTransaction(ormService, func() error {
		flusher := ormService.NewFlusher()

		roleEntity := &entity.RoleEntity{
			Name:      request.Name,
			CreatedAt: now,
		}

		flusher.Track(roleEntity)

		if err := createPrivileges(flusher, roleEntity, request.Resources, resourcesMapping, permissionsMapping, now); err != nil {
			return err
		}

		return flusher.FlushWithCheck()
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateRole(c *gin.Context, roleID *acl.RoleRequestDTO, request *acl.CreateOrUpdateRoleRequestDTO) error {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	roleEntity := &entity.RoleEntity{}
	if !ormService.LoadByID(roleID.ID, roleEntity) {
		return fmt.Errorf("role with ID: %d not found", roleID.ID)
	}

	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
	if err != nil {
		return err
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterUint("RoleID", roleEntity.ID)

	privilegeEntitiesToDelete := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(query, beeorm.NewPager(1, 1000), &privilegeEntitiesToDelete)

	now := service.DI().Clock().Now()

	err = helper.DBTransaction(ormService, func() error {
		flusher := ormService.NewFlusher()

		for _, privilegeEntity := range privilegeEntitiesToDelete {
			flusher.Delete(privilegeEntity)
		}

		if err := flusher.FlushWithCheck(); err != nil {
			return err
		}

		flusher = ormService.NewFlusher()

		roleEntity.Name = request.Name

		flusher.Track(roleEntity)

		if err := createPrivileges(flusher, roleEntity, request.Resources, resourcesMapping, permissionsMapping, now); err != nil {
			return err
		}

		return flusher.FlushWithCheck()
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteRole(c *gin.Context, roleID *acl.RoleRequestDTO) error {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	roleEntity := &entity.RoleEntity{}
	if !ormService.LoadByID(roleID.ID, roleEntity) {
		return fmt.Errorf("role with ID: %d not found", roleID.ID)
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterUint("RoleID", roleEntity.ID)

	privilegeEntitiesToDelete := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(query, beeorm.NewPager(1, 1000), &privilegeEntitiesToDelete)

	err := helper.DBTransaction(ormService, func() error {
		flusher := ormService.NewFlusher()

		flusher.Delete(roleEntity)

		for _, privilegeEntity := range privilegeEntitiesToDelete {
			flusher.Delete(privilegeEntity)
		}

		return flusher.FlushWithCheck()
	})
	if err != nil {
		return err
	}

	return nil
}

func PostAssignRoleToUserAction(c *gin.Context, getUserFunc func() beeorm.Entity, request *acl.AssignRoleToUserRequestDTO) error {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	roleEntity := &entity.RoleEntity{}
	if !ormService.LoadByID(request.RoleID, roleEntity) {
		return fmt.Errorf("role with ID: %d not found", request.RoleID)
	}

	userEntity := getUserFunc()
	if !ormService.LoadByID(request.UserID, userEntity) {
		return fmt.Errorf("user with ID: %d not found", request.UserID)
	}

	userWithSettableRole, ok := userEntity.(UserRoleSetter)
	if !ok {
		panic("user entity does not implement UserRoleSetter interface")
	}

	userWithSettableRole.SetRole(roleEntity)

	userEntityWithNewRole, _ := userWithSettableRole.(beeorm.Entity)

	ormService.Flush(userEntityWithNewRole)

	return nil
}

type resourceMapping map[uint64]*entity.ResourceEntity

type permissionMapping map[uint64]*entity.PermissionEntity

// nolint // info
func validateResourcesAndPermissions(ormService *datalayer.ORM, resources []*acl.RoleResourceRequestDTO) (resourceMapping, permissionMapping, error) {
	resourceIDs := make([]uint64, len(resources))
	permissionIDs := make([]uint64, 0)

	for i, resource := range resources {
		resourceIDs[i] = resource.ResourceID

		permissionIDs = append(permissionIDs, resource.PermissionIDs...)
	}

	resourcesQuery := redisearch.NewRedisSearchQuery()
	resourcesQuery.FilterUint("ID", resourceIDs...)

	resourceEntities := make([]*entity.ResourceEntity, 0)
	ormService.RedisSearch(resourcesQuery, beeorm.NewPager(1, 1000), &resourceEntities)

	if len(resourceEntities) != len(resourceIDs) {
		return nil, nil, fmt.Errorf("some of the provided resources is not found")
	}

	permissionsQuery := redisearch.NewRedisSearchQuery()
	permissionsQuery.FilterUint("ID", permissionIDs...)

	permissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(permissionsQuery, beeorm.NewPager(1, 4000), &permissionEntities)

	if len(permissionEntities) != len(permissionIDs) {
		return nil, nil, fmt.Errorf("some of the provided permissions is not found")
	}

	resourcesMapping := resourceMapping{}

	for _, resourceEntity := range resourceEntities {
		resourcesMapping[resourceEntity.ID] = resourceEntity
	}

	permissionsMapping := permissionMapping{}

	for _, permissionEntity := range permissionEntities {
		permissionsMapping[permissionEntity.ID] = permissionEntity
	}

	return resourcesMapping, permissionsMapping, nil
}

func createPrivileges(
	flusher beeorm.Flusher,
	roleEntity *entity.RoleEntity,
	resources []*acl.RoleResourceRequestDTO,
	resourcesMapping resourceMapping,
	permissionsMapping permissionMapping,
	now time.Time,
) error {
	for _, resource := range resources {
		if len(resource.PermissionIDs) == 0 {
			continue
		}

		resourceEntity, ok := resourcesMapping[resource.ResourceID]
		if !ok {
			return fmt.Errorf("resource with ID: %d not found in mapping", resource.ResourceID)
		}

		privilegeEntity := &entity.PrivilegeEntity{
			RoleID:     roleEntity,
			ResourceID: resourceEntity,
			CreatedAt:  now,
		}

		for _, permissionID := range resource.PermissionIDs {
			permissionEntity, ok := permissionsMapping[permissionID]
			if !ok {
				return fmt.Errorf("permission with ID: %d not found in mapping", permissionID)
			}

			if permissionEntity.ResourceID.ID != resourceEntity.ID {
				return fmt.Errorf("permission with ID: %d does not belong to resource with ID: %d", permissionID, resource.ResourceID)
			}

			privilegeEntity.PermissionIDs = append(privilegeEntity.PermissionIDs, permissionEntity)
		}

		flusher.Track(privilegeEntity)
	}

	return nil
}
