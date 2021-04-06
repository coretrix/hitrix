package helper

import (
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GraphqlErr(path []string, msg string) error {
	pathContext := graphql.NewPathWithField(strings.Join(path, "."))
	return &gqlerror.Error{
		Message: msg,
		Path:    pathContext.Path(),
	}
}
