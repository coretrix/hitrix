package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf("%d", i))
	})
}

func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("string failed to be parsed: %v", err)
		}
		return uint64(i), nil
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case json.Number:
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, fmt.Errorf("json.Number failed to be parsed: %v", err)
		}
		return uint64(i), nil
	default:
		return 0, fmt.Errorf("%T is not an int", v)
	}
}
