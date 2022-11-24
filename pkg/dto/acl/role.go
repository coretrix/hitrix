package acl

type RolesResponseDTO struct {
	Roles []*RoleResponseDTO
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
