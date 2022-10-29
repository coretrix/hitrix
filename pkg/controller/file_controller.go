package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/file"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	fileModel "github.com/coretrix/hitrix/pkg/model/file"
	"github.com/coretrix/hitrix/pkg/response"
)

type FileController struct {
}

// @Description Upload
// @Tags File
// @Accept mpfd
// @Param body formData file.RequestDTOUploadFile true "Request in formData"
// @Param file formData file true "The file"
// @Router /file/upload/ [post]
// @Security BearerAuth
// @Success 200 {object} file.File
// @Failure 400 {object} response.Error
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Something bad happened"
func (controller *FileController) PostUploadImageAction(c *gin.Context) {
	req := &file.RequestDTOUploadFile{}

	err := binding.ShouldBind(c, req)

	if errorhandling.HandleError(c, err) {
		return
	}

	request, closeFile, err := req.ToUploadImage()

	if errorhandling.HandleError(c, err) {
		return
	}

	defer func() {
		_ = closeFile()
	}()

	res, err := fileModel.CreateFile(c.Request.Context(), request)

	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, res)
}
