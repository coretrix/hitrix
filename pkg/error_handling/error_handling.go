package errorhandling

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/errors"
	"github.com/coretrix/hitrix/pkg/response"
)

func HandleError(c *gin.Context, err error) bool {
	errType, ok := err.(errors.FieldErrors)
	if ok && errType != nil {
		response.ErrorResponseFields(c, errType, nil)

		return true
	}

	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return true
	}

	return false
}
