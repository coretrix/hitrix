package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/errors"
)

const ResponseBody = "response_body"

type Error struct {
	GlobalError string            `json:"GlobalError,omitempty"`
	FieldsError map[string]string `json:"FieldsError,omitempty"`
	Result      *interface{}      `json:"Result,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.Set(ResponseBody, data)

	if data != nil {
		c.JSON(http.StatusOK, data)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func NotFoundResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{})
}

func ErrorResponseGlobal(c *gin.Context, globalError interface{}, data interface{}) {
	result := &Error{}

	if data != nil {
		result.Result = &data
	}

	err, ok := globalError.(*errors.PermissionError)
	if ok {
		c.Set(ResponseBody, err.Error())

		c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
		return
	}

	err1, ok1 := globalError.(error)
	if ok1 {
		result.GlobalError = err1.Error()
	} else {
		result.GlobalError = globalError.(string)
	}

	c.Set(ResponseBody, result)

	c.JSON(http.StatusBadRequest, result)
}

func ErrorResponseFields(c *gin.Context, fieldsError errors.FieldErrors, data interface{}) {
	result := &Error{
		FieldsError: fieldsError,
	}

	if data != nil {
		result.Result = &data
	}

	c.Set(ResponseBody, result)

	c.JSON(http.StatusBadRequest, result)
}
