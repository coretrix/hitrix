package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service"
)

type ReadinessController struct {
}

func (controller *ReadinessController) GetReadinessAction(c *gin.Context) {
	ormService := service.DI().OrmEngine()

	var res int8

	has := ormService.GetMysql().QueryRow(beeorm.NewWhere("SELECT 1"), &res)
	if !has || res != 1 {
		c.JSON(503, gin.H{"error": "mysql do not respond"})

		return
	}

	ormService.GetRedis().Set("ping", 1, helper.Minute)

	_, has = ormService.GetRedis().Get("ping")
	if !has {
		c.JSON(503, gin.H{"error": "redis do not respond"})

		return
	}

	c.JSON(200, "OK")
}

func (controller *ReadinessController) GetLivenessAction(c *gin.Context) {
	c.JSON(200, "OK")
}
