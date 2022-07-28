package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/view/account"
)

func AuthorizeDevUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := account.IsValidDevToken(c, c.Request.Header.Get("Authorization"))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
	}
}

func AuthorizeWithDevRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := account.IsValidDevRefreshToken(c, c.Request.Header.Get("Authorization"))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
	}
}
