package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
)

type ClockworkController struct {
}

func (controller *ClockworkController) GetIndexAction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NotFoundResponse(c)
		return
	}
	profilerService := service.DI().ClockWorkForContext(c.Request.Context())

	c.JSON(http.StatusOK, profilerService.GetSavedData(id))
}
