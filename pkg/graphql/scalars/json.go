package scalars

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalJSON(json json.RawMessage) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, string(json))
	})
}

func UnmarshalJSON(v interface{}) (json.RawMessage, error) {
	switch v := v.(type) {
	case json.RawMessage:
		return v, nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("%T is not json", v)
	}
}
