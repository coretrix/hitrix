package acl

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
)

func RolesColumns() []crud.Column {
	return []crud.Column{
		{
			Key:        "ID",
			Label:      "ID",
			Searchable: false,
			Sortable:   false,
			Visible:    true,
		},
		{
			Key:        "Name",
			Label:      "Name",
			Searchable: false,
			Sortable:   false,
			Visible:    true,
		},
	}
}

func ListRoles(c *gin.Context, request *crud.ListRequest) *acl.RolesResponseDTO {
	cols := RolesColumns()
	crudService := service.DI().Crud()

	searchParams := crudService.ExtractListParams(cols, request)

	query := crudService.GenerateListRedisSearchQuery(searchParams)
	if len(searchParams.Sort) == 0 {
		query.Sort("ID", false)
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	allRoleEntities := make([]*entity.RoleEntity, 0)
	total := ormService.RedisSearch(&allRoleEntities, query, beeorm.NewPager(searchParams.Page, searchParams.PageSize))

	result := &acl.RolesResponseDTO{
		Total:   int(total),
		Columns: cols,
	}

	for _, roleEntity := range allRoleEntities {
		result.Rows = append(result.Rows, &acl.RoleResponseDTO{
			ID:   roleEntity.ID,
			Name: roleEntity.Name,
		})
	}

	return result
}

func GetRole(c *gin.Context, request *acl.RoleRequestDTO) (*acl.RoleResponseDTO, error) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	query.FilterUint("RoleID", request.ID)

	allPrivilegeEntities := make([]*entity.PrivilegeEntity, 0)
	ormService.RedisSearch(&allPrivilegeEntities, query, beeorm.NewPager(1, 4000), "RoleID", "ResourceID", "PermissionIDs")

	if len(allPrivilegeEntities) == 0 {
		return nil, fmt.Errorf("role with ID: %d not found", request.ID)
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

	return result, nil
}
