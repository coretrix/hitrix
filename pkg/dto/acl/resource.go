package acl

type ResourcesResponseDTO struct {
	Resources []*ResourceResponseDTO
}

type ResourceResponseDTO struct {
	ID          uint64
	Name        string
	Permissions []*PermissionResponseDTO
}
