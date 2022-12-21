package acl

import "github.com/coretrix/hitrix/service/component/crud"

type RolesResponseDTO struct {
	Rows    []*RoleResponseDTO
	Total   int
	Columns []crud.Column
}

type RoleResponseDTO struct {
	ID        uint64
	Name      string
	Resources []*ResourceResponseDTO `json:",omitempty"`
}

type RoleRequestDTO struct {
	ID uint64 `binding:"required"`
}

type CreateOrUpdateRoleRequestDTO struct {
	Name      string                    `binding:"required"`
	Resources []*RoleResourceRequestDTO `binding:"required,min=1,dive"`
}

type RoleResourceRequestDTO struct {
	ResourceID  uint64                      `binding:"required"`
	Permissions []*RolePermissionRequestDTO `binding:"required,min=1,dive"`
}

type RolePermissionRequestDTO struct {
	PermissionID uint64 `binding:"required"`
}

type AssignRoleToUserRequestDTO struct {
	UserID uint64 `binding:"required"`
	RoleID uint64 `binding:"required"`
}
