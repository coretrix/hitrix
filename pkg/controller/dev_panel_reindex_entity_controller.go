package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/dto/indexes"
	"github.com/coretrix/hitrix/pkg/model/redissearch"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
)

func (controller *DevPanelController) PostRedisSearchReindexEntity(c *gin.Context) {
	request := indexes.RequestDTOReindexEntity{}
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ErrorResponseGlobal(c, "ID must be a positive integer encoded as a decimal string", nil)

		return
	}

	id, err := strconv.ParseUint(request.ID, 10, 64)
	if err != nil || id == 0 || strconv.FormatUint(id, 10) != request.ID {
		response.ErrorResponseGlobal(c, "ID must be a positive integer up to 18446744073709551615", nil)

		return
	}

	appService := service.DI().App()
	pool := c.Param("pool")
	allowedPool := false
	if appService.DevPanel != nil && appService.RedisPools != nil {
		for _, searchPool := range appService.RedisPools.Search {
			if searchPool == pool {
				allowedPool = true
				break
			}
		}
	}
	if !allowedPool {
		response.ErrorResponseGlobal(c, "unknown Redis Search pool", nil)

		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())
	found, err := redissearch.ReindexEntity(ormService, pool, c.Param("index"), id)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}
	if !found {
		c.JSON(http.StatusNotFound, response.Error{GlobalError: "record not found or excluded from this index"})

		return
	}

	response.SuccessResponse(c, gin.H{"ID": request.ID})
}
