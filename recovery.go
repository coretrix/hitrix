package hitrix

import (
	"bytes"
	"github.com/coretrix/hitrix/service"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorLogger := service.DI().ErrorLogger()
				requestBody := c.Request.Context().Value(service.RequestBodyKey).([]byte)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
				errorLogger.LogErrorWithRequest(c, r)

			}
		}()

		c.Next()
	}
}
