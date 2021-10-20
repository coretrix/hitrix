package binding

import (
	"fmt"

	"github.com/coretrix/hitrix/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ShouldBindJSON(c *gin.Context, form interface{}) error {
	if c.Request.Body == nil {
		return fmt.Errorf("body cannot be nil")
	}

	if err := c.ShouldBindBodyWith(form, binding.JSON); err != nil {
		res := errors.HandleErrors(err)
		if res != nil {
			return res
		}

		return err
	}

	return nil
}

func ShouldBindQuery(c *gin.Context, form interface{}) error {
	if err := c.ShouldBindQuery(form); err != nil {
		res := errors.HandleErrors(err)
		if res != nil {
			return res
		}

		return err
	}

	return nil
}
