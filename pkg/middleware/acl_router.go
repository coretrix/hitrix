package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/controller"
)

func ACLRouter(ginEngine *gin.Engine) {
	v1Group := ginEngine.Group("/v1/")

	var aclController *controller.ACLController
	{
		aclGroup := v1Group.Group("/acl")

		aclGroup.GET("/resources/", aclController.ListResourcesAction)
		aclGroup.GET("/role/:ID/", aclController.GetRoleAction)
		aclGroup.POST("/roles/", aclController.ListRolesAction)
		aclGroup.POST("/role/", aclController.CreateRoleAction)
		aclGroup.PUT("/role/:ID/", aclController.UpdateRoleAction)
		aclGroup.DELETE("/role/:ID/", aclController.DeleteRoleAction)
		aclGroup.POST("/assign-role/", aclController.PostAssignRoleToUserAction(func() beeorm.Entity {
			return &entity.AdminUserEntity{}
		}))
	}
}
