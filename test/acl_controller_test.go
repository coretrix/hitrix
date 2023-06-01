package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xorcare/pointer"

	entityExample "github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	aclView "github.com/coretrix/hitrix/pkg/view/acl"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock/mocks"
	"github.com/coretrix/hitrix/service/component/crud"
	registryMocks "github.com/coretrix/hitrix/service/registry/mocks"
)

func TestListResourcesAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	resource1 := CreateResource(flusher, map[string]interface{}{})
	flusher.Flush()
	resource2 := CreateResource(flusher, map[string]interface{}{"Name": "car"})

	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource1, "Name": "create"})
	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource1, "Name": "view"})

	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource2, "Name": "unlock"})
	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource2, "Name": "lock"})
	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource2, "Name": "drive"})

	flusher.Flush()

	got := &acl.ResourcesResponseDTO{}

	err := SendHTTPRequest(ctx, http.MethodGet, "/acl/resources/", false, got)
	assert.Nil(t, err)

	want := &acl.ResourcesResponseDTO{
		Resources: []*acl.ResourceResponseDTO{
			{
				ID:   1,
				Name: "user",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   1,
						Name: "create",
					},
					{
						ID:   2,
						Name: "view",
					},
				},
			},
			{
				ID:   2,
				Name: "car",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   3,
						Name: "unlock",
					},
					{
						ID:   4,
						Name: "lock",
					},
					{
						ID:   5,
						Name: "drive",
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestListRolesAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	resource1 := CreateResource(flusher, map[string]interface{}{})

	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource1, "Name": "create"})
	CreatePermission(flusher, map[string]interface{}{"ResourceID": resource1, "Name": "view"})

	CreateRole(flusher, map[string]interface{}{})
	CreateRole(flusher, map[string]interface{}{"Name": "super-admin"})
	CreateRole(flusher, map[string]interface{}{"Name": "super-mega-admin"})

	flusher.Flush()

	got := &acl.RolesResponseDTO{}

	request := &crud.ListRequest{
		Page:     pointer.Int(1),
		PageSize: pointer.Int(2),
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/roles/", request, false, got)
	assert.Nil(t, err)

	want := &acl.RolesResponseDTO{
		Total:   3,
		Columns: aclView.RolesColumns(),
		Rows: []*acl.RoleResponseDTO{
			{
				ID:   1,
				Name: "admin",
			},
			{
				ID:   2,
				Name: "super-admin",
			},
		},
		PageContext: &acl.ResourcesResponseDTO{
			Resources: []*acl.ResourceResponseDTO{
				{
					ID:   1,
					Name: "user",
					Permissions: []*acl.PermissionResponseDTO{
						{
							ID:   1,
							Name: "create",
						},
						{
							ID:   2,
							Name: "view",
						},
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestGetRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	role := CreateRole(flusher, map[string]interface{}{})

	resource := CreateResource(flusher, map[string]interface{}{})

	permission1 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(flusher, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	flusher.Flush()

	got := &acl.RoleResponseDTO{}

	err := SendHTTPRequest(ctx, http.MethodGet, "/acl/role/1/", false, got)
	assert.Nil(t, err)

	want := &acl.RoleResponseDTO{
		ID:   1,
		Name: "admin",
		Resources: []*acl.ResourceResponseDTO{
			{
				ID:   1,
				Name: "user",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   1,
						Name: "create",
					},
					{
						ID:   2,
						Name: "view",
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestCreateRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	resource := CreateResource(flusher, map[string]interface{}{})

	permission1 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	flusher.Flush()

	request := &acl.CreateOrUpdateRoleRequestDTO{
		Name: "admin",
		Resources: []*acl.RoleResourceRequestDTO{
			{
				ResourceID:    resource.ID,
				PermissionIDs: []uint64{permission1.ID, permission2.ID},
			},
		},
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/role/", request, false, nil)
	assert.Nil(t, err)

	privilegeEntity := &entity.PrivilegeEntity{}
	assert.True(t, ormService.LoadByID(1, privilegeEntity))

	assert.Equal(t, privilegeEntity.RoleID.ID, uint64(1))
	assert.Equal(t, privilegeEntity.ResourceID.ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[0].ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[1].ID, uint64(2))

	fakeClock.AssertExpectations(t)
}

func TestUpdateRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	role := CreateRole(flusher, map[string]interface{}{})

	resource := CreateResource(flusher, map[string]interface{}{})

	permission1 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(flusher, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	flusher.Flush()

	request := &acl.CreateOrUpdateRoleRequestDTO{
		Name: "super-admin",
		Resources: []*acl.RoleResourceRequestDTO{
			{
				ResourceID:    resource.ID,
				PermissionIDs: []uint64{permission2.ID},
			},
		},
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPut, "/acl/role/1/", request, false, nil)
	assert.Nil(t, err)

	privilegeEntity := &entity.PrivilegeEntity{}
	assert.False(t, ormService.LoadByID(1, privilegeEntity))
	privilegeEntity = &entity.PrivilegeEntity{}
	assert.True(t, ormService.LoadByID(2, privilegeEntity, "RoleID"))

	assert.Equal(t, privilegeEntity.RoleID.ID, uint64(1))
	assert.Equal(t, privilegeEntity.RoleID.Name, "super-admin")
	assert.Equal(t, privilegeEntity.ResourceID.ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[0].ID, permission2.ID)
	assert.Equal(t, len(privilegeEntity.PermissionIDs), 1)

	fakeClock.AssertExpectations(t)
}

func TestDeleteRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	role := CreateRole(flusher, map[string]interface{}{})

	resource := CreateResource(flusher, map[string]interface{}{})

	permission1 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(flusher, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(flusher, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	flusher.Flush()

	err := SendHTTPRequest(ctx, http.MethodDelete, "/acl/role/1/", false, nil)
	assert.Nil(t, err)

	roleEntity := &entity.RoleEntity{}
	assert.True(t, ormService.LoadByID(1, roleEntity))
	assert.Equal(t, roleEntity.FakeDelete, true)

	privilegeEntity := &entity.PrivilegeEntity{}
	assert.True(t, ormService.LoadByID(1, privilegeEntity))
	assert.Equal(t, privilegeEntity.FakeDelete, true)

	fakeClock.AssertExpectations(t)
}

func TestPostAssignRoleToUserAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", nil, mockServices, nil)

	ormService := service.DI().OrmEngine().Clone()
	flusher := ormService.NewFlusher()

	role1 := CreateRole(flusher, map[string]interface{}{})
	role2 := CreateRole(flusher, map[string]interface{}{"Name": "super-admin"})

	user := CreateAdminUser(flusher, map[string]interface{}{"RoleID": role1})

	flusher.Flush()

	request := &acl.AssignRoleToUserRequestDTO{
		UserID: user.ID,
		RoleID: role2.ID,
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/assign-role/", request, false, nil)
	assert.Nil(t, err)

	userEntity := &entityExample.AdminUserEntity{}
	assert.True(t, ormService.LoadByID(1, userEntity, "RoleID"))
	assert.Equal(t, userEntity.RoleID.ID, role2.ID)
	assert.Equal(t, userEntity.RoleID.Name, role2.Name)

	fakeClock.AssertExpectations(t)
}
