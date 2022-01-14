package binding

import (
	"github.com/gin-gonic/gin/binding"

	"github.com/coretrix/hitrix/pkg/errors"
)

func ValidateStruct(s interface{}) error {
	err := binding.Validator.ValidateStruct(s)
	if err != nil {
		res := errors.HandleErrors(err)
		if res != nil {
			return res
		}

		return err
	}

	return nil
}
