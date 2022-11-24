package controller

import (
	"fmt"
	"time"

	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
)

type ACLController struct {
}

// @Description List resources
// @Tags ACL
// @Router /acl/resources/ [get]
// @Success 200 {object} acl.ResourcesResponseDTO
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) ListResourcesAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	allPermissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(&allPermissionEntities, beeorm.NewRedisSearchQuery(), beeorm.NewPager(1, 4000), "ResourceID")

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

	result := &acl.ResourcesResponseDTO{}

	for _, resourceDTO := range resourceDTOsMapping {
		result.Resources = append(result.Resources, resourceDTO)
	}

	response.SuccessResponse(c, result)
}

// @Description List roles
// @Tags ACL
// @Router /acl/roles/ [get]
// @Success 200 {object} acl.RolesResponseDTO
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) ListRolesAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	allRoleEntities := make([]*entity.RoleEntity, 0)
	ormService.RedisSearch(&allRoleEntities, beeorm.NewRedisSearchQuery(), beeorm.NewPager(1, 4000))

	result := &acl.RolesResponseDTO{}

	for _, roleEntity := range allRoleEntities {
		result.Roles = append(result.Roles, &acl.RoleResponseDTO{
			ID:   roleEntity.ID,
			Name: roleEntity.Name,
		})
	}

	response.SuccessResponse(c, result)
}

// @Description Get role
// @Tags ACL
// @Param ID path string true "ID"
// @Router /acl/role/{ID}/ [get]
// @Success 200 {object} acl.RoleResponseDTO
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) GetRoleAction(c *gin.Context) {
	request := &acl.RoleRequestDTO{}

	err := binding.ShouldBindURI(c, request)
	if errorhandling.HandleError(c, err) {
		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	query.FilterUint("RoleID", request.ID)

	allPrivilegeEntities := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(&allPrivilegeEntities, query, beeorm.NewPager(1, 4000), "RoleID", "ResourceID", "PermissionIDs")

	if len(allPrivilegeEntities) == 0 {
		errorhandling.HandleError(c, fmt.Errorf("role with ID: %d not found", request.ID))

		return
	}

	result := &acl.RoleResponseDTO{}

	for i, privilegeEntity := range allPrivilegeEntities {
		if i == 0 {
			result.ID = privilegeEntity.RoleID.ID
			result.Name = privilegeEntity.RoleID.Name
		}

		resource := &acl.ResourceResponseDTO{
			ID:   privilegeEntity.ResourceID.ID,
			Name: privilegeEntity.ResourceID.Name,
		}

		for _, permission := range privilegeEntity.PermissionIDs {
			resource.Permissions = append(resource.Permissions, &acl.PermissionResponseDTO{
				ID:   permission.ID,
				Name: permission.Name,
			})
		}

		result.Resources = append(result.Resources, resource)
	}

	response.SuccessResponse(c, result)
}

// @Description Create role
// @Tags ACL
// @Param body body acl.CreateOrUpdateRoleRequestDTO true "Request in body"
// @Router /acl/role/ [post]
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) CreateRoleAction(c *gin.Context) {
	request := &acl.CreateOrUpdateRoleRequestDTO{}

	err := binding.ShouldBindJSON(c, request)
	if errorhandling.HandleError(c, err) {
		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
	if errorhandling.HandleError(c, err) {
		return
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
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, nil)
}

// @Description Update role
// @Tags ACL
// @Param ID path string true "ID"
// @Param body body acl.CreateOrUpdateRoleRequestDTO true "Request in body"
// @Router /acl/role/{ID}/ [put]
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) UpdateRoleAction(c *gin.Context) {
	roleID := &acl.RoleRequestDTO{}

	err := binding.ShouldBindURI(c, roleID)
	if errorhandling.HandleError(c, err) {
		return
	}

	request := &acl.CreateOrUpdateRoleRequestDTO{}

	err = binding.ShouldBindJSON(c, request)
	if errorhandling.HandleError(c, err) {
		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	roleEntity := &entity.RoleEntity{}
	if !ormService.LoadByID(roleID.ID, roleEntity) {
		errorhandling.HandleError(c, fmt.Errorf("role with ID: %d not found", roleID.ID))

		return
	}

	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
	if errorhandling.HandleError(c, err) {
		return
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterUint("RoleID", roleEntity.ID)

	privilegeEntitiesToDelete := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(&privilegeEntitiesToDelete, query, beeorm.NewPager(1, 1000))

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
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, nil)
}

// @Description Delete role
// @Tags ACL
// @Param ID path string true "ID"
// @Router /acl/role/{ID}/ [delete]
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) DeleteRoleAction(c *gin.Context) {
	roleID := &acl.RoleRequestDTO{}

	err := binding.ShouldBindURI(c, roleID)
	if errorhandling.HandleError(c, err) {
		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	roleEntity := &entity.RoleEntity{}
	if !ormService.LoadByID(roleID.ID, roleEntity) {
		errorhandling.HandleError(c, err)

		return
	}

	query := beeorm.NewRedisSearchQuery()
	query.FilterUint("RoleID", roleEntity.ID)

	privilegeEntitiesToDelete := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(&privilegeEntitiesToDelete, query, beeorm.NewPager(1, 1000))

	err = helper.DBTransaction(ormService, func() error {
		flusher := ormService.NewFlusher()

		flusher.Delete(roleEntity)

		for _, privilegeEntity := range privilegeEntitiesToDelete {
			flusher.Delete(privilegeEntity)
		}

		return flusher.FlushWithCheck()
	})
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, nil)
}

type resourceDTOsMapping map[uint64]*acl.ResourceResponseDTO

type resourceMapping map[uint64]*entity.ResourceEntity

type permissionMapping map[uint64]*entity.PermissionEntity

func validateResourcesAndPermissions(ormService *beeorm.Engine, resources []*acl.RoleResourceRequestDTO) (resourceMapping, permissionMapping, error) {
	resourceIDs := make([]uint64, len(resources))
	permissionIDs := make([]uint64, 0)

	for i, resource := range resources {
		resourceIDs[i] = resource.ResourceID

		for _, permission := range resource.Permissions {
			permissionIDs = append(permissionIDs, permission.PermissionID)
		}
	}

	resourcesQuery := beeorm.NewRedisSearchQuery()
	resourcesQuery.FilterUint("ID", resourceIDs...)

	resourceEntities := make([]*entity.ResourceEntity, 0)
	ormService.RedisSearch(&resourceEntities, resourcesQuery, beeorm.NewPager(1, 1000))

	if len(resourceEntities) != len(resourceIDs) {
		return nil, nil, fmt.Errorf("some of the provided resources is not found")
	}

	permissionsQuery := beeorm.NewRedisSearchQuery()
	permissionsQuery.FilterUint("ID", permissionIDs...)

	permissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(&permissionEntities, permissionsQuery, beeorm.NewPager(1, 4000))

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
		resourceEntity, ok := resourcesMapping[resource.ResourceID]
		if !ok {
			return fmt.Errorf("resource with ID: %d not found in mapping", resource.ResourceID)
		}

		privilegeEntity := &entity.PrivilegeEntity{
			RoleID:     roleEntity,
			ResourceID: resourceEntity,
			CreatedAt:  now,
		}

		for _, permission := range resource.Permissions {
			permissionEntity, ok := permissionsMapping[permission.PermissionID]
			if !ok {
				return fmt.Errorf("permission with ID: %d not found in mapping", permission.PermissionID)
			}

			if permissionEntity.ResourceID.ID != resourceEntity.ID {
				return fmt.Errorf("permission with ID: %d does not belong to resource with ID: %d", permission.PermissionID, resource.ResourceID)
			}

			privilegeEntity.PermissionIDs = append(privilegeEntity.PermissionIDs, permissionEntity)
		}

		flusher.Track(privilegeEntity)
	}

	return nil
}
