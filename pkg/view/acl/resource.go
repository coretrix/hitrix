package acl

import (
	"sort"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func ListResources(c *gin.Context) *acl.ResourcesResponseDTO {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := redisearch.NewRedisSearchQuery()
	query.Sort("ID", false)

	allPermissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearchMany(query, beeorm.NewPager(1, 4000), &allPermissionEntities, "ResourceID")

	resourceDTOsMapping := resourceDTOsMapping{}

	for _, permissionEntity := range allPermissionEntities {
		dto, ok := resourceDTOsMapping[permissionEntity.ResourceID.ID]
		if !ok {
			dto = &acl.ResourceResponseDTO{
				ID:          permissionEntity.ResourceID.ID,
				Name:        permissionEntity.ResourceID.Name,
				Permissions: make([]*acl.PermissionResponseDTO, 0),
			}
		}

		dto.Permissions = append(dto.Permissions, &acl.PermissionResponseDTO{
			ID:   permissionEntity.ID,
			Name: permissionEntity.Name,
		})

		resourceDTOsMapping[permissionEntity.ResourceID.ID] = dto
	}

	resultDTOs := make([]*acl.ResourceResponseDTO, 0)

	for _, resourceDTO := range resourceDTOsMapping {
		resultDTOs = append(resultDTOs, resourceDTO)
	}

	sort.Slice(resultDTOs, func(i, j int) bool {
		return resultDTOs[i].ID < resultDTOs[j].ID
	})

	return &acl.ResourcesResponseDTO{Resources: resultDTOs}
}

type UserRoleGetter interface {
	GetRole() *entity.RoleEntity
}

func ListUserResources(c *gin.Context, getUserFunc func(c *gin.Context) beeorm.Entity) *acl.ResourcesResponseDTO {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	userEntity := getUserFunc(c)

	userWithGettableRole, ok := userEntity.(UserRoleGetter)
	if !ok {
		panic("user entity does not implement UserRoleSetter interface")
	}

	query := redisearch.NewRedisSearchQuery()
	query.FilterUint("RoleID", userWithGettableRole.GetRole().ID)

	privilegeEntities := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearchMany(query, beeorm.NewPager(1, 4000), &privilegeEntities, "ResourceID")

	for _, privilege := range privilegeEntities {
		permissionIDs := make([]uint64, len(privilege.PermissionIDs))

		for i, permission := range privilege.PermissionIDs {
			permissionIDs[i] = permission.ID
		}

		permissionEntities := make([]*entity.PermissionEntity, 0)
		ormService.LoadByIDs(permissionIDs, &permissionEntities)

		privilege.PermissionIDs = permissionEntities
	}

	resourceDTOsMapping := resourceDTOsMapping{}

	for _, privilegeEntity := range privilegeEntities {
		dto, ok := resourceDTOsMapping[privilegeEntity.ResourceID.ID]
		if !ok {
			dto = &acl.ResourceResponseDTO{
				ID:          privilegeEntity.ResourceID.ID,
				Name:        privilegeEntity.ResourceID.Name,
				Permissions: make([]*acl.PermissionResponseDTO, 0),
			}
		}

		for _, permissionEntity := range privilegeEntity.PermissionIDs {
			dto.Permissions = append(dto.Permissions, &acl.PermissionResponseDTO{
				ID:   permissionEntity.ID,
				Name: permissionEntity.Name,
			})
		}

		resourceDTOsMapping[privilegeEntity.ResourceID.ID] = dto
	}

	resultDTOs := make([]*acl.ResourceResponseDTO, 0)

	for _, resourceDTO := range resourceDTOsMapping {
		resultDTOs = append(resultDTOs, resourceDTO)
	}

	sort.Slice(resultDTOs, func(i, j int) bool {
		return resultDTOs[i].ID < resultDTOs[j].ID
	})

	return &acl.ResourcesResponseDTO{Resources: resultDTOs}
}

type resourceDTOsMapping map[uint64]*acl.ResourceResponseDTO
