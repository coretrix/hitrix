package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/acl"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	aclModel "github.com/coretrix/hitrix/pkg/model/acl"
	"github.com/coretrix/hitrix/pkg/response"
	aclView "github.com/coretrix/hitrix/pkg/view/acl"
	"github.com/coretrix/hitrix/service/component/crud"
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
	response.SuccessResponse(c, aclView.ListResources(c))
}

// @Description List roles
// @Tags ACL
// @Param body body crud.ListRequest true "Request in body"
// @Router /acl/roles/ [post]
// @Success 200 {object} acl.RolesResponseDTO
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) ListRolesAction(c *gin.Context) {
	request := &crud.ListRequest{}

	err := binding.ShouldBindJSON(c, request)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, aclView.ListRoles(c, request))
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

	result, err := aclView.GetRole(c, request)
	if errorhandling.HandleError(c, err) {
		return
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

	err = aclModel.CreateRole(c, request)
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

	err = aclModel.UpdateRole(c, roleID, request)
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

	err = aclModel.DeleteRole(c, roleID)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, nil)
}

// @Description Assign role to user
// @Tags ACL
// @Param body body acl.AssignRoleToUserRequestDTO true "Request in body"
// @Router /acl/assign-role/ [post]
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
// @Security BearerAuth
func (controller *ACLController) PostAssignRoleToUserAction(getUserFunc func() beeorm.Entity) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &acl.AssignRoleToUserRequestDTO{}

		err := binding.ShouldBindJSON(c, request)
		if errorhandling.HandleError(c, err) {
			return
		}

		err = aclModel.PostAssignRoleToUserAction(c, getUserFunc, request)
		if errorhandling.HandleError(c, err) {
			return
		}

		response.SuccessResponse(c, nil)
	}
}
