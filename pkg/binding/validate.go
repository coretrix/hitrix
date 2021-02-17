package binding

import (
	errors "github.com/coretrix/hitrix/pkg/error"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ValidateStruct(c *gin.Context, s interface{}) error {
	err := binding.Validator.ValidateStruct(s)
	if err != nil {
		res := errors.HandleErrors(err, c)
		if res != nil {
			return res
		}

		return err
	}

	return nil
}
