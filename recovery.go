package hitrix

import (
	"log"

	"github.com/gin-gonic/gin"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorLogger, has := DIC().ErrorLogger()
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
