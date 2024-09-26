package hitrix

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.Status(http.StatusInternalServerError)

				service.DI().ErrorLogger().LogPanicWithRequest(c, r)
			}
		}()

		c.Next()
	}
}
