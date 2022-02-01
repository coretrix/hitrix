package gql

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/coretrix/hitrix/service/component/localize"
)

const Language = "language"

type IGQLInterface interface {
	GraphqlErrPath(cxt context.Context, path []string, msg string, args ...interface{}) error
	GraphqlErr(cxt context.Context, msg string, args ...interface{}) error
}

type gqlService struct {
	localizeService localize.ILocalizer
}

func (t gqlService) GraphqlErrPath(cxt context.Context, path []string, msg string, args ...interface{}) error {
	pathContext := graphql.NewPathWithField(strings.Join(path, "."))

	err := t.setError(cxt, msg, args)
	err.Path = pathContext.Path()

	return err
}

func (t gqlService) GraphqlErr(cxt context.Context, msg string, args ...interface{}) error {
	return t.setError(cxt, msg, args)
}

func (t gqlService) setError(cxt context.Context, msg string, args ...interface{}) *gqlerror.Error {
	err := &gqlerror.Error{}

	if t.localizeService != nil && cxt.Value(Language) != nil {
		translatedMessage := t.localizeService.T(cxt.Value(Language).(string), msg)
		err.Message = fmt.Sprintf(translatedMessage, args...)
	} else {
		err.Message = fmt.Sprintf(msg, args...)
	}

	err.Extensions = map[string]interface{}{
		"code": msg,
	}

	return err
}

func NewGqlService(localizeService localize.ILocalizer) IGQLInterface {
	return &gqlService{localizeService: localizeService}
}
