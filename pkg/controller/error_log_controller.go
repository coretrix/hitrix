package controller

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/pkg/response"

	"github.com/gin-gonic/gin"
)

type ErrorLogController struct {
}

func (controller *ErrorLogController) GetErrors(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c.Request.Context())
	if !has {
		panic("orm is not registered")
	}

	type errorRow struct {
		File    string
		Line    int
		AppName string
		Message string
		Stack   string
		Counter int
		Time    string
	}

	data := ormService.GetRedis().HGetAll(hitrix.GroupError)

	errorsList := map[string]*errorRow{}

	for key, value := range data {
		// TODO: fix this hack
		if len(value) == 0 {
			continue
		}

		splitKeys := strings.Split(key, ":")

		if _, ok := errorsList[splitKeys[0]]; !ok {
			errorsList[splitKeys[0]] = &errorRow{}
		}

		if len(splitKeys) == 1 {
			errorMessage := &hitrix.ErrorMessage{}
			err := json.Unmarshal([]byte(value), errorMessage)
			if err != nil {
				panic(err)
			}
			errorsList[splitKeys[0]].Stack = string(errorMessage.Stack)
			errorsList[splitKeys[0]].File = errorMessage.File
			errorsList[splitKeys[0]].Message = errorMessage.Message
			errorsList[splitKeys[0]].Line = errorMessage.Line
			errorsList[splitKeys[0]].AppName = errorMessage.AppName
		} else if len(splitKeys) == 2 {
			if splitKeys[1] == "time" {
				i, _ := strconv.ParseInt(value, 10, 64)
				errorsList[splitKeys[0]].Time = time.Unix(i, 0).String()
			} else if splitKeys[1] == "counter" {
				counter, _ := strconv.Atoi(value)
				errorsList[splitKeys[0]].Counter = counter
			}
		}
	}

	response.SuccessResponse(c, errorsList)
}

func (controller *ErrorLogController) DeleteError(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c.Request.Context())
	if !has {
		panic("orm is not registered")
	}

	id := c.Param("id")
	if len(id) <= 0 {
		response.ErrorResponseGlobal(c, "missing id", nil)
		return
	}
	ormService.GetRedis().HDel(hitrix.GroupError, id)
	ormService.GetRedis().HDel(hitrix.GroupError, id+":time")
	ormService.GetRedis().HDel(hitrix.GroupError, id+":counter")

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) DeleteAllErrors(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c.Request.Context())
	if !has {
		panic("orm is not registered")
	}

	ormService.GetRedis().Del(hitrix.GroupError)

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) Panic(_ *gin.Context) {
	panic("Forced Panic")
}
