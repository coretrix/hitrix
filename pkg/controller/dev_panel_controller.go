package controller

import (
	"fmt"
	"sort"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix/service"

	errors "github.com/coretrix/hitrix/pkg/error"
	accountModel "github.com/coretrix/hitrix/pkg/model/account"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/pkg/view/account"

	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm/tools"
)

type MenuItem struct {
	Label string
	URL   string
	Icon  string
}

type DevPanelController struct {
}

//func (controller *DevPanelController) GetActionListAction(c *gin.Context) {
//	actions := []*MenuItem{
//		{
//			Label: "Clear Cache",
//			URL:   "/dev/clear-cache/",
//			Icon:  "mdiCached",
//		},
//	}
//
//	c.JSON(200, actions)
//}

func (controller *DevPanelController) CreateAdminUserAction(c *gin.Context) {
	passwordService, has := service.DI().Password()
	if !has {
		panic("password is not registered")
	}

	ormService, has := service.DI().OrmEngine()
	if !has {
		panic("password is not registered")
	}

	username := c.Query("username")
	pass := c.Query("password")

	if username == "" || pass == "" {
		response.ErrorResponseGlobal(c, "username and password query parameters are required", nil)
		return
	}
	adminEntity := service.DI().App().DevPanel.UserEntity

	passwordHash, err := passwordService.HashPassword(pass)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
	}

	adminTableSchema := ormService.GetRegistry().GetTableSchemaForEntity(adminEntity)
	response.SuccessResponse(c, fmt.Sprintf(`INSERT INTO %s (Email, Password) VALUES("%s", "%s")`, adminTableSchema.GetTableName(), username, passwordHash))
}

func (controller *DevPanelController) PostLoginDevPanelAction(c *gin.Context) {
	fmt.Println("PostLoginDevPanelAction")
	loginForm := accountModel.LoginDevForm{}
	token, refreshToken, err := loginForm.Login(c)

	errType, ok := err.(errors.FieldErrors)

	if ok && errType != nil {
		response.ErrorResponseFields(c, errType, nil)
		return
	}

	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) PostGenerateTokenAction(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())
	if !has {
		panic("orm is not registered")
	}

	devPanelUserEntity := c.MustGet(account.LoggedDevPanelUserEntity).(app.DevPanelUserEntity)

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) GetClearCacheAction(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	redisService := ormService.GetRedis()

	redisService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetClearRedisStreamsAction(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	devPanel := service.DI().App().DevPanel
	if devPanel == nil || devPanel.PoolStream == nil {
		panic("stream pool is not defined")
	}

	redisStreamsService := ormService.GetRedis(*devPanel.PoolStream)
	redisStreamsService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) DeleteRedisStreamAction(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	devPanel := service.DI().App().DevPanel
	if devPanel == nil || devPanel.PoolStream == nil {
		panic("stream pool is not defined")
	}

	redisStreamService := ormService.GetRedis(*devPanel.PoolStream)

	name := c.Param("name")
	if name == "" {
		panic("provide stream name")
	}

	redisStreamService.XTrim(name, 0, false)

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetAlters(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	dbService := ormService.GetMysql()

	alters := ormService.GetAlters()

	result := make([]string, len(alters))
	force := c.Query("force")
	if force != "" {
		redisService := ormService.GetRedis()
		redisService.FlushDB()
	}

	for i, alter := range alters {
		if force != "" {
			dbService.Exec(alter.SQL)
		} else {
			result[i] = alter.SQL
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisStreams(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	stats := tools.GetRedisStreamsStatistics(ormService)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Stream < stats[j].Stream
	})
	response.SuccessResponse(c, stats)
}

func (controller *DevPanelController) GetRedisStatistics(c *gin.Context) {
	ormService, has := service.DI().OrmEngineForContext(c.Request.Context())

	if !has {
		panic("orm is not registered")
	}

	stats := tools.GetRedisStatistics(ormService)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].RedisPool < stats[j].RedisPool
	})
	response.SuccessResponse(c, stats)
}
