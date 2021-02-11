package response

import (
	"net/http"

	errors "github.com/coretrix/hitrix/pkg/error"
	"github.com/gin-gonic/gin"
)

type Error struct {
	GlobalError string            `json:"GlobalError,omitempty"`
	FieldsError map[string]string `json:"FieldsError,omitempty"`
	Result      *interface{}      `json:"Result,omitempty"`
}

func SuccessAPIResponse(c *gin.Context, data interface{}) {
	if data != nil {
		c.JSON(http.StatusOK, data)

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"Result": data,
		})

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
		c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
		return
	}

	err1, ok1 := globalError.(error)
	if ok1 {
		c.JSON(http.StatusBadRequest, err1.Error())
		return
	}

	result.GlobalError = globalError.(string)

	c.JSON(http.StatusBadRequest, result)
}

func ErrorResponseFields(c *gin.Context, fieldsError errors.FieldErrors, data interface{}) {
	result := &Error{
		FieldsError: fieldsError,
	}

	if data != nil {
		result.Result = &data
	}

	c.JSON(http.StatusBadRequest, result)
}
