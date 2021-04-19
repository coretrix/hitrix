package hitrix

import (
	"context"
	"reflect"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ValidateDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver, rules string) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rules string) (interface{}, error) {
		originalValue, err := next(ctx)
		if err != nil {
			return nil, err
		}

		switch v := originalValue.(type) {
		case string:
			if v == "" {
				return v, nil
			}
		default:
			if v == nil || reflect.ValueOf(v).IsNil() {
				return v, nil
			}
		}

		errs := helper.NewValidator().Validate(originalValue, rules)

		for _, e := range errs {
			graphql.AddError(ctx, &gqlerror.Error{
				Path:    graphql.GetPath(ctx),
				Message: "Field" + e.Error(),
			})
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
