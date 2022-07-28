package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

func AuthorizeWithQueryParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretFromQueryString := c.Query("Secret")

		secretFromConfig, ok := service.DI().Config().String("query_param.secret")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		if secretFromConfig != secretFromQueryString {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Next()
	}
}
