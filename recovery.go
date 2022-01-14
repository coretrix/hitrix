package hitrix

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.Request.Body = ioutil.NopCloser(
					bytes.NewReader(c.Request.Context().Value(service.RequestBodyKey).([]byte)))

				service.DI().ErrorLogger().LogErrorWithRequest(c, r)
			}
		}()

		c.Next()
	}
}
