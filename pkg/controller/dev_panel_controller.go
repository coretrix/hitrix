package controller

import (
	"sort"

	"github.com/coretrix/hitrix"
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

func (controller *DevPanelController) PostLoginDevPanelAction(c *gin.Context) {
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
	ormService, has := hitrix.DIC().OrmEngineForContext(c)
	if !has {
		panic("orm is not registered")
	}

	devPanelUserEntity := c.MustGet(account.LoggedDevPanelUserEntity).(hitrix.DevPanelUserEntity)

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
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	redisService := ormService.GetRedis()

	redisService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetClearRedisStreamsAction(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	devPanel := hitrix.DIC().App().DevPanel()
	if devPanel == nil || devPanel.PoolStream == nil {
		panic("stream pool is not defined")
	}

	redisStreamsService := ormService.GetRedis(*devPanel.PoolStream)
	redisStreamsService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) DeleteRedisStreamAction(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	devPanel := hitrix.DIC().App().DevPanel()
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
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	dbService := ormService.GetMysql()

	alters := ormService.GetAlters()

	result := make([]string, 0)
	force := c.Query("force")
	for _, alter := range alters {
		if force != "" {
			redisService := ormService.GetRedis()

			redisService.FlushDB()

			dbService.Exec(alter.SQL)
		} else {
			result = append(result, alter.SQL)
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisStreams(c *gin.Context) {
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

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
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	stats := tools.GetRedisStatistics(ormService)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].RedisPool < stats[j].RedisPool
	})
	response.SuccessResponse(c, stats)
}
