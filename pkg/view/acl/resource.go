package acl

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func ListResources(c *gin.Context) *acl.ResourcesResponseDTO {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := beeorm.NewRedisSearchQuery()
	query.Sort("ID", false)

	allPermissionEntities := make([]*entity.PermissionEntity, 0)
	ormService.RedisSearch(&allPermissionEntities, query, beeorm.NewPager(1, 4000), "ResourceID")

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

type resourceDTOsMapping map[uint64]*acl.ResourceResponseDTO
