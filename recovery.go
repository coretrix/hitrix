package hitrix

import (
	"log"

	"github.com/coretrix/hitrix/service"

	"github.com/gin-gonic/gin"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorLogger, has := service.DI().ErrorLogger()
				if has {
					errorLogger.LogRecover(r)
				} else {
					log.Println(r.(string))
				}
			}
		}()

		c.Next()
	}
}
