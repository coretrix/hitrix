package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/list"
	translationDTO "github.com/coretrix/hitrix/pkg/dto/translation"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	model "github.com/coretrix/hitrix/pkg/model/translation"
	"github.com/coretrix/hitrix/pkg/response"
	view "github.com/coretrix/hitrix/pkg/view/translation"
)

type TranslationController struct {
}

// @Description Translation List
// @Description	Parameters:
// @Description	Page     *int `binding:"required"`
// @Description	PageSize *int `binding:"required"`
// @Description	Search   map[string]interface{}
// @Description	SearchOR map[string]interface{}
// @Description	Sort     map[string]interface{}
// @Tags Translation
// @Param body body crud.ListRequest true "Request in body"
// @Router /translation/list/ [post]
// @Security BearerAuth
// @Success 200 {object} translation.ResponseDTOList
// @Failure 400 {object} response.Error
// @Failure 500 "Something bad happened"
func (controller *TranslationController) PostTranslationListAction(c *gin.Context) {
	request := list.RequestDTOList{}

	err := binding.ShouldBindJSON(c, &request)
	if errorhandling.HandleError(c, err) {
		return
	}

	res, err := view.List(c.Request.Context(), request)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, res)
}

// @Description Create Translation
// @Tags Translation
// @Router /translation/create/ [post]
// @Param body body translation.RequestCreateTranslation true "Request in body"
// @Security BearerAuth
// @Success 200 {object} translation.ResponseTranslation
// @Failure 400 {object} response.Error
// @Failure 500 "Something bad happened"
func (controller *TranslationController) PostCreateTranslationAction(c *gin.Context) {
	request := &translationDTO.RequestCreateTranslation{}
	err := binding.ShouldBindJSON(c, request)

	if errorhandling.HandleError(c, err) {
		return
	}

	data, err := model.Create(c.Request.Context(), request)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, data)
}

// @Description Update Translation
// @Tags Translation
// @Param body body translation.RequestUpdateTranslation true "Request in body"
// @Router /translation/update/{ID}/ [post]
// @Param ID path string true "id"
// @Security BearerAuth
// @Success 200 {object} translation.ResponseTranslation
// @Failure 400 {object} response.Error
// @Failure 500 "Something bad happened"
func (controller *TranslationController) PostUpdateTranslationAction(c *gin.Context) {
	request := &translationDTO.RequestUpdateTranslation{}
	err := binding.ShouldBindJSON(c, request)

	if errorhandling.HandleError(c, err) {
		return
	}

	requestTranslationID := &translationDTO.RequestDTOTranslationID{}

	err = binding.ShouldBindURI(c, requestTranslationID)
	if errorhandling.HandleError(c, err) {
		return
	}

	data, err := model.Update(c.Request.Context(), request, requestTranslationID.ID)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, data)
}

// @Description Delete Translation
// @Tags Translation
// @Router /translation/delete/{ID}/ [delete]
// @Param ID path string true "id"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 500 "Something bad happened"
func (controller *TranslationController) DeleteTranslationAction(c *gin.Context) {
	requestTranslationID := &translationDTO.RequestDTOTranslationID{}

	err := binding.ShouldBindURI(c, requestTranslationID)
	if errorhandling.HandleError(c, err) {
		return
	}

	err = model.Delete(c.Request.Context(), requestTranslationID.ID)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, nil)
}

// @Description Get Translation
// @Tags Translation
// @Router /translation/{ID}/ [get]
// @Param ID path string true "id"
// @Security BearerAuth
// @Success 200 {object} translation.ResponseTranslation
// @Failure 400 {object} response.Error
// @Failure 500 "Something bad happened"
func (controller *TranslationController) GetTranslationAction(c *gin.Context) {
	requestTranslationID := &translationDTO.RequestDTOTranslationID{}

	err := binding.ShouldBindURI(c, requestTranslationID)
	if errorhandling.HandleError(c, err) {
		return
	}

	data, err := view.Get(c.Request.Context(), requestTranslationID.ID)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, data)
}
