package hitrix

import (
	"context"
	"fmt"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vEn "github.com/go-playground/validator/v10/translations/en"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ValidateDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver, rules string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rules string) (interface{}, error) {
		originalValue, err := next(ctx)
		if err != nil {
			return nil, err
		}

		if originalValue == nil || reflect.ValueOf(originalValue).IsNil() {
			return originalValue, nil
		}

		validate := validator.New()

		english := en.New()
		uni := ut.New(english, english)
		trans, _ := uni.GetTranslator("en")
		_ = vEn.RegisterDefaultTranslations(validate, trans)

		err = validate.Var(originalValue, rules)
		errs := translateError(err, trans)

		if err != nil {
			for _, e := range errs {
				graphql.AddError(ctx, &gqlerror.Error{
					Path:    graphql.GetPath(ctx),
					Message: "Field" + e.Error(),
				})
			}

			return originalValue, nil
		}

		return originalValue, nil
	}
}

func Validate(ctx context.Context, callback func() bool) bool {
	if graphql.GetErrors(ctx) != nil {
		return false
	}

	if callback != nil {
		return callback()
	}

	return true
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
