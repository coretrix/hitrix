package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldErrors map[string]string

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "unauthorized"
}

type PermissionError struct {
	Message string
}

func (e *PermissionError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "permission denied"
}

func HandleErrors(formErrors interface{}, c *gin.Context) error {
	fErrors, ok := formErrors.(validator.ValidationErrors)

	if !ok {
		return nil
	}

	var fe FieldErrors = make(map[string]string)
	for _, errors := range fErrors {
		fe[errors.Field()] = errors.Translate(nil)
	}

	return fe
}

func HandleCustomErrors(formErrors map[string]string, c *gin.Context) error {
	var fe FieldErrors = make(map[string]string)
	for field, msg := range formErrors {
		fe[field] = msg
	}

	return fe
}

func (fe FieldErrors) Error() string {
	var result string
	for _, val := range fe {
		result += val + "\n\r"
	}
	return result
}
