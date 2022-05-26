package binding

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/leebenson/conform"

	"github.com/coretrix/hitrix/pkg/errors"
)

func ValidateStruct(s interface{}) error {
	if err := conform.Strings(s); err != nil {
		return err
	}

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
