package controller

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ErrorLogController struct {
}

type errorRow struct {
	File    string
	Line    int
	AppName string
	Request string
	Message string
	Stack   string
	Counter int
	Time    string
}

func (controller *ErrorLogController) getByGroup(c *gin.Context, group string) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())
	data := ormService.GetRedis().HGetAll(group)

	errorsList := map[string]*errorRow{}

	for key, value := range data {
		// TODO: fix this hack
		if len(value) == 0 {
			continue
		}

		entryID := key
		field := ""
		lastSplitIndex := strings.LastIndex(key, ":")
		if lastSplitIndex > 0 {
			entryID = key[0:lastSplitIndex]
			field = key[lastSplitIndex+1:]
		}

		if _, ok := errorsList[entryID]; !ok {
			errorsList[entryID] = &errorRow{}
		}

		if field == "" {
			errorMessage := &errorlogger.ErrorMessage{}

			err := json.Unmarshal([]byte(value), errorMessage)
			if err != nil {
				panic(err)
			}

			errorsList[entryID].Request = string(errorMessage.Request)
			errorsList[entryID].Stack = string(errorMessage.Stack)
			errorsList[entryID].File = errorMessage.File
			errorsList[entryID].Message = errorMessage.Message
			errorsList[entryID].Line = errorMessage.Line
			errorsList[entryID].AppName = errorMessage.AppName
		} else if field == "time" {
			i, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
			errorsList[entryID].Time = time.Unix(i, 0).String()
		} else if field == "counter" {
			counter, _ := strconv.Atoi(strings.TrimSpace(value))
			errorsList[entryID].Counter = counter
		}
	}

	response.SuccessResponse(c, errorsList)
}

func (controller *ErrorLogController) deleteSingleByGroup(c *gin.Context, group string) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	id := c.Param("id")
	if len(id) <= 0 {
		response.ErrorResponseGlobal(c, "missing id", nil)

		return
	}

	ormService.GetRedis().HDel(group, id)
	ormService.GetRedis().HDel(group, id+":time")
	ormService.GetRedis().HDel(group, id+":counter")

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) deleteAllByGroup(c *gin.Context, group string) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())
	ormService.GetRedis().Del(group)

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) GetErrors(c *gin.Context) {
	controller.getByGroup(c, errorlogger.GroupError)
}

func (controller *ErrorLogController) GetCounters(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())
	countByGroup := func(group string) int {
		data := ormService.GetRedis().HGetAll(group)
		counterEntries := 0
		messageEntries := 0

		for key, value := range data {
			if len(value) == 0 {
				continue
			}

			lastSplitIndex := strings.LastIndex(key, ":")
			if lastSplitIndex <= 0 {
				messageEntries++

				continue
			}

			field := key[lastSplitIndex+1:]
			if field == "counter" {
				counterEntries++
			}
		}

		if counterEntries > 0 {
			return counterEntries
		}

		return messageEntries
	}

	response.SuccessResponse(c, map[string]int{
		"errors":              countByGroup(errorlogger.GroupError),
		"warnings":            countByGroup(errorlogger.GroupWarning),
		"missingTranslations": countByGroup(errorlogger.GroupMissingTranslation),
	})
}

func (controller *ErrorLogController) DeleteError(c *gin.Context) {
	controller.deleteSingleByGroup(c, errorlogger.GroupError)
}

func (controller *ErrorLogController) DeleteAllErrors(c *gin.Context) {
	controller.deleteAllByGroup(c, errorlogger.GroupError)
}

func (controller *ErrorLogController) GetWarnings(c *gin.Context) {
	controller.getByGroup(c, errorlogger.GroupWarning)
}

func (controller *ErrorLogController) DeleteWarning(c *gin.Context) {
	controller.deleteSingleByGroup(c, errorlogger.GroupWarning)
}

func (controller *ErrorLogController) DeleteAllWarnings(c *gin.Context) {
	controller.deleteAllByGroup(c, errorlogger.GroupWarning)
}

func (controller *ErrorLogController) GetMissingTranslations(c *gin.Context) {
	controller.getByGroup(c, errorlogger.GroupMissingTranslation)
}

func (controller *ErrorLogController) DeleteMissingTranslation(c *gin.Context) {
	controller.deleteSingleByGroup(c, errorlogger.GroupMissingTranslation)
}

func (controller *ErrorLogController) DeleteAllMissingTranslations(c *gin.Context) {
	controller.deleteAllByGroup(c, errorlogger.GroupMissingTranslation)
}

func (controller *ErrorLogController) Panic(_ *gin.Context) {
	panic("Forced Panic")
}
