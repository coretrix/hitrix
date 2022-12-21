# ACL

You can use ACL feature by including 4 hitrix entities in you orm init:

```go
type RoleEntity struct {
    beeorm.ORM `orm:"table=roles;redisCache;redisSearch=search_pool"`
    ID         uint64    `orm:"sortable"`
    Name       string    `orm:"required;searchable;unique=Name_FakeDelete:1"`
    CreatedAt  time.Time `orm:"time=true"`
    FakeDelete bool      `orm:"unique=Name_FakeDelete:2"`
}

type ResourceEntity struct {
    beeorm.ORM `orm:"table=resources;redisCache;redisSearch=search_pool"`
    ID         uint64    `orm:"searchable"`
    Name       string    `orm:"required;searchable;unique=Name_FakeDelete:1"`
    CreatedAt  time.Time `orm:"time=true"`
    FakeDelete bool      `orm:"unique=Name_FakeDelete:2"`
}

type PermissionEntity struct {
    beeorm.ORM `orm:"table=permissions;redisCache;redisSearch=search_pool"`
    ID         uint64          `orm:"searchable;sortable"`
    ResourceID *ResourceEntity `orm:"required;searchable;unique=ResourceID_Name_FakeDelete:1"`
    Name       string          `orm:"required;searchable;unique=ResourceID_Name_FakeDelete:2"`
    CreatedAt  time.Time       `orm:"time=true"`
    FakeDelete bool            `orm:"unique=ResourceID_Name_FakeDelete:3"`
}

type PrivilegeEntity struct {
    beeorm.ORM    `orm:"table=privileges;redisCache;redisSearch=search_pool"`
    ID            uint64
    RoleID        *RoleEntity         `orm:"required;searchable;unique=RoleID_ResourceID_FakeDelete:1"`
    ResourceID    *ResourceEntity     `orm:"required;searchable;unique=RoleID_ResourceID_FakeDelete:2"`
    PermissionIDs []*PermissionEntity `orm:"required;searchable"`
    CreatedAt     time.Time           `orm:"time=true"`
    FakeDelete    bool                `orm:"unique=RoleID_ResourceID_FakeDelete:3"`
}
```
After you do this, in you local (to the specific project) user entity, you can include foreign key to
role entity like this:

```go

type UserEntity struct {
    beeorm.ORM `orm:"table=users;log=log_db_pool;redisCache;redisSearch=search_pool"`
    ID         uint64
    RoleID     *hitrixEntity.RoleEntity `orm:"required"`
}
```

Now you are ready to use the ACL feature by using these endpoints:

```go
var aclController *hitrixController.ACLController
{
	v1Group.GET("/acl/resources/", aclController.ListResourcesAction)
	v1Group.GET("/acl/role/:ID/", aclController.GetRoleAction)
	v1Group.POST("/acl/roles/", aclController.ListRolesAction)
	v1Group.POST("/acl/role/", aclController.CreateRoleAction)
	v1Group.PUT("/acl/role/:ID/", aclController.UpdateRoleAction)
	v1Group.DELETE("/acl/role/:ID/", aclController.DeleteRoleAction)
	v1Group.POST("/acl/assign-role/", aclController.PostAssignRoleToUserAction(func() beeorm.Entity {
		return &entity.AdminUserEntity{}
	}))
}
```

These endpoints allow you to configure different combinations of roles, resources and permissions.
After you configure your roles and permissions, you can start using the ACL by using middleware:

```go
func ACL(resource string, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminUserEntity, ok := ioc.GetAdminUserService().GetSession(c.Request.Context())
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		ormService := service.DI().OrmEngineForContext(c.Request.Context())

		if !acl.ACL(ormService, adminUserEntity.User.RoleID, resource, permissions...) {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Next()
	}
}

```

For example if you want to make the roles and permissions endpoints protected by ACL, you can modify
the router to look like this:

```go
var aclController *hitrixController.ACLController
aclGroup := v1Group.Group("/acl/").Use(authMiddleware.AuthorizeWithHeaderStrict())
{
	aclGroup.GET("resources/", authMiddleware.ACL(constant.ResourceResource, constant.PermissionView), aclController.ListResourcesAction)
	aclGroup.GET("role/:ID/", authMiddleware.ACL(constant.ResourceRole, constant.PermissionView), aclController.GetRoleAction)
	aclGroup.POST("roles/", authMiddleware.ACL(constant.ResourceRole, constant.PermissionView), aclController.ListRolesAction)
	aclGroup.POST("role/", authMiddleware.ACL(constant.ResourceRole, constant.PermissionModify), aclController.CreateRoleAction)
	aclGroup.PUT("role/:ID/", authMiddleware.ACL(constant.ResourceRole, constant.PermissionModify), aclController.UpdateRoleAction)
	aclGroup.DELETE("role/:ID/", authMiddleware.ACL(constant.ResourceRole, constant.PermissionModify), aclController.DeleteRoleAction)
	aclGroup.POST("assign-role/", authMiddleware.ACL(constant.ResourceAdminUser, constant.PermissionAssignRole), aclController.PostAssignRoleToUserAction(func() beeorm.Entity {
		return &entity.AdminUserEntity{}
	}))
}
```

You can see how this middleware is used on each endpoint, and it checks if the logged user role has
the required privileges that are specified on each endpoint.
