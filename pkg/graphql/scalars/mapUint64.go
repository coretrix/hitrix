package scalars

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalMapUint64(m map[uint64]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		marshaled, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		_, err = w.Write(marshaled)
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalMapUint64(v interface{}) (map[uint64]interface{}, error) {
	switch v := v.(type) {
	case map[uint64]interface{}:
		return v, nil
	case json.RawMessage:
		m := map[uint64]interface{}{}
		if err := json.Unmarshal(v, &m); err != nil {
			return nil, err
		}
		return m, nil
	case string:
		m := map[uint64]interface{}{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("%T is not a map", v)
	}
}
