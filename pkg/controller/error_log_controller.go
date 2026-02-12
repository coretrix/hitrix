package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ErrorLogController struct {
}

type errorRow struct {
	File       string
	Line       int
	AppName    string
	Request    string
	Message    string
	Stack      string
	Counter    int
	Time       string
	TicketLink string
}

func (controller *ErrorLogController) getErrorLogRedis(c *gin.Context) *beeorm.RedisCache {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())
	appService := service.DI().App()

	if appService.RedisPools != nil && appService.RedisPools.Persistent != "" {
		return ormService.GetRedis(appService.RedisPools.Persistent)
	}

	return ormService.GetRedis()
}

func (controller *ErrorLogController) getByGroup(c *gin.Context, group string) {
	data := controller.getErrorLogRedis(c).HGetAll(group)

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
		} else if field == "ticket_link" {
			errorsList[entryID].TicketLink = value
		}
	}

	response.SuccessResponse(c, errorsList)
}

func (controller *ErrorLogController) deleteSingleByGroup(c *gin.Context, group string) {
	redisStorage := controller.getErrorLogRedis(c)

	id := c.Param("id")
	if len(id) <= 0 {
		response.ErrorResponseGlobal(c, "missing id", nil)

		return
	}

	redisStorage.HDel(group, id)
	redisStorage.HDel(group, id+":time")
	redisStorage.HDel(group, id+":counter")
	redisStorage.HDel(group, id+":ticket_link")

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) deleteAllByGroup(c *gin.Context, group string) {
	controller.getErrorLogRedis(c).Del(group)

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) GetErrors(c *gin.Context) {
	controller.getByGroup(c, errorlogger.GroupError)
}

func (controller *ErrorLogController) GetCounters(c *gin.Context) {
	redisStorage := controller.getErrorLogRedis(c)
	countByGroup := func(group string) int {
		data := redisStorage.HGetAll(group)
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

func (controller *ErrorLogController) CreateErrorJiraTicket(c *gin.Context) {
	controller.createJiraTicketByGroup(c, errorlogger.GroupError, "Error")
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

func (controller *ErrorLogController) CreateWarningJiraTicket(c *gin.Context) {
	controller.createJiraTicketByGroup(c, errorlogger.GroupWarning, "Warning")
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

func (controller *ErrorLogController) CreateMissingTranslationJiraTicket(c *gin.Context) {
	controller.createJiraTicketByGroup(c, errorlogger.GroupMissingTranslation, "Missing translation")
}

func (controller *ErrorLogController) DeleteAllMissingTranslations(c *gin.Context) {
	controller.deleteAllByGroup(c, errorlogger.GroupMissingTranslation)
}

func (controller *ErrorLogController) Panic(_ *gin.Context) {
	panic("Forced Panic")
}

func (controller *ErrorLogController) createJiraTicketByGroup(c *gin.Context, group, itemType string) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponseGlobal(c, "missing id", nil)

		return
	}

	redisStorage := controller.getErrorLogRedis(c)
	if existingTicket, has := redisStorage.HGet(group, id+":ticket_link"); has && existingTicket != "" {
		response.SuccessResponse(c, map[string]string{
			"TicketLink": existingTicket,
		})

		return
	}

	rawMessage, has := redisStorage.HGet(group, id)
	if !has || rawMessage == "" {
		response.ErrorResponseGlobal(c, "error log item is missing", nil)

		return
	}

	errorMessage := &errorlogger.ErrorMessage{}
	if err := json.Unmarshal([]byte(rawMessage), errorMessage); err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	jiraURL, has := service.DI().Config().String("jira.url")
	if !has || jiraURL == "" {
		response.ErrorResponseGlobal(c, "missing jira.url config", nil)

		return
	}

	jiraToken, has := service.DI().Config().String("jira.token")
	if !has || jiraToken == "" {
		response.ErrorResponseGlobal(c, "missing jira.token config", nil)

		return
	}

	jiraUser, has := service.DI().Config().String("jira.user")
	if !has || jiraUser == "" {
		response.ErrorResponseGlobal(c, "missing jira.user config", nil)

		return
	}

	projectKey, has := service.DI().Config().String("jira.project_key")
	if !has || projectKey == "" {
		response.ErrorResponseGlobal(c, "missing jira.project_key config", nil)

		return
	}

	issueType, has := service.DI().Config().String("jira.issue_type")
	if !has || issueType == "" {
		response.ErrorResponseGlobal(c, "missing jira.issue_type config", nil)

		return
	}

	timeValue := ""
	if unixTime, hasTime := redisStorage.HGet(group, id+":time"); hasTime && unixTime != "" {
		timestamp, err := strconv.ParseInt(strings.TrimSpace(unixTime), 10, 64)
		if err == nil {
			timeValue = time.Unix(timestamp, 0).Format(time.RFC3339)
		}
	}
	if timeValue == "" {
		timeValue = time.Now().Format(time.RFC3339)
	}

	description := fmt.Sprintf(
		"Type: %s\n\n"+
			"File: %s\n"+
			"Line: %d\n"+
			"AppName: %s\n"+
			"Time: %s\n\n"+
			"Message:\n%s\n\n"+
			"Request:\n%s\n\n"+
			"Stack:\n%s",
		itemType,
		errorMessage.File,
		errorMessage.Line,
		errorMessage.AppName,
		timeValue,
		errorMessage.Message,
		string(errorMessage.Request),
		string(errorMessage.Stack),
	)

	requestBody := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": projectKey,
			},
			"summary":     truncateString(fmt.Sprintf("[%s] %s", itemType, errorMessage.Message), 255),
			"description": truncateString(description, 15000),
			"issuetype": map[string]string{
				"name": issueType,
			},
		},
	}

	if assigneeAccountID, has := service.DI().Config().String("jira.assignee_account_id"); has && assigneeAccountID != "" {
		requestBody["fields"].(map[string]interface{})["assignee"] = map[string]string{
			"accountId": assigneeAccountID,
		}
	}

	bodyBinary, err := json.Marshal(requestBody)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(jiraURL, "/")+"/rest/api/2/issue", bytes.NewBuffer(bodyBinary))
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(jiraUser, jiraToken)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		response.ErrorResponseGlobal(
			c,
			fmt.Sprintf("jira create issue failed with status %d: %s", resp.StatusCode, truncateString(string(respBody), 1000)),
			nil,
		)

		return
	}

	type issueResponse struct {
		Key string `json:"key"`
	}

	issue := &issueResponse{}
	if err := json.Unmarshal(respBody, issue); err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	ticketLink := strings.TrimRight(jiraURL, "/") + "/browse/" + issue.Key
	redisStorage.HSet(group, id+":ticket_link", ticketLink)

	response.SuccessResponse(c, map[string]string{
		"TicketLink": ticketLink,
		"TicketKey":  issue.Key,
	})
}

func truncateString(value string, limit int) string {
	if len(value) <= limit {
		return value
	}

	return value[:limit]
}
