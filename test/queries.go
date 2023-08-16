package main

import (
	"time"

	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/example/entity"
	entityHitrix "github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func CreateAdminUser(flusher beeorm.Flusher, row map[string]interface{}) *entity.AdminUserEntity {
	userEntity := &entity.AdminUserEntity{}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "RoleID":
				userEntity.RoleID = value.(*entityHitrix.RoleEntity)
			}
		}
	}

	flusher.Track(userEntity)

	return userEntity
}

func CreateRole(flusher beeorm.Flusher, row map[string]interface{}) *entityHitrix.RoleEntity {
	roleEntity := &entityHitrix.RoleEntity{
		Name:      "admin",
		CreatedAt: service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "Name":
				roleEntity.Name = value.(string)
			case "CreatedAt":
				roleEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	flusher.Track(roleEntity)

	return roleEntity
}

func CreateResource(flusher beeorm.Flusher, row map[string]interface{}) *entityHitrix.ResourceEntity {
	resourceEntity := &entityHitrix.ResourceEntity{
		Name:      "user",
		CreatedAt: service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "Name":
				resourceEntity.Name = value.(string)
			case "CreatedAt":
				resourceEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	flusher.Track(resourceEntity)

	return resourceEntity
}

func CreatePermission(flusher beeorm.Flusher, row map[string]interface{}) *entityHitrix.PermissionEntity {
	permissionEntity := &entityHitrix.PermissionEntity{
		ResourceID: nil,
		Name:       "view",
		CreatedAt:  service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "ResourceID":
				permissionEntity.ResourceID = value.(*entityHitrix.ResourceEntity)
			case "Name":
				permissionEntity.Name = value.(string)
			case "CreatedAt":
				permissionEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	flusher.Track(permissionEntity)

	return permissionEntity
}

func CreatePrivilege(flusher beeorm.Flusher, row map[string]interface{}) {
	permissionEntity := &entityHitrix.PrivilegeEntity{
		RoleID:        nil,
		ResourceID:    nil,
		PermissionIDs: nil,
		CreatedAt:     service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "RoleID":
				permissionEntity.RoleID = value.(*entityHitrix.RoleEntity)
			case "ResourceID":
				permissionEntity.ResourceID = value.(*entityHitrix.ResourceEntity)
			case "PermissionIDs":
				permissionEntity.PermissionIDs = value.([]*entityHitrix.PermissionEntity)
			case "CreatedAt":
				permissionEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	flusher.Track(permissionEntity)
}
