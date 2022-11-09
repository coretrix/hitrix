package hitrix

import (
	"context"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/coretrix/hitrix/pkg/binding"
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

		errs := binding.NewValidator().Validate(originalValue, rules)
		if len(errs) > 0 {
			if len(errs) > 1 {
				for i := 1; i < len(errs); i++ {
					graphql.AddError(ctx, &gqlerror.Error{
						Path:    graphql.GetPath(ctx),
						Message: "Field" + errs[i].Error(),
					})
				}
			}

			return nil, &gqlerror.Error{
				Path:    graphql.GetPath(ctx),
				Message: "Field" + errs[0].Error(),
			}
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
