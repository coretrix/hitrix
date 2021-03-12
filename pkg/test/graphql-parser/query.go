package graphqlparser

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"strings"
)

func constructQuery(v interface{}, variables map[string]interface{}) string {
	query := query(v)
	if len(variables) > 0 {
		return "query(" + queryArguments(variables) + ")" + query
	}
	return query
}

func constructMutation(v interface{}, variables map[string]interface{}) string {
	query := query(v)
	if len(variables) > 0 {
		return "mutation(" + queryArguments(variables) + ")" + query
	}
	return "mutation" + query
}

func queryArguments(variables map[string]interface{}) string {
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		// nolint
		io.WriteString(&buf, "$")
		// nolint
		io.WriteString(&buf, k)
		// nolint
		io.WriteString(&buf, ":")
		writeArgumentType(&buf, reflect.TypeOf(variables[k]), true)
	}
	return buf.String()
}

func writeArgumentType(w io.Writer, t reflect.Type, value bool) {
	if t.Kind() == reflect.Ptr {
		writeArgumentType(w, t.Elem(), false)
		return
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		// nolint
		io.WriteString(w, "[")
		writeArgumentType(w, t.Elem(), true)
		// nolint
		io.WriteString(w, "]")
	default:
		// nolint
		io.WriteString(w, strings.Title(t.Name()))
	}

	if value {
		// nolint
		io.WriteString(w, "!")
	}
}

func query(v interface{}) string {
	var buf bytes.Buffer
	writeQuery(&buf, reflect.TypeOf(v), false)
	return buf.String()
}

func writeQuery(w io.Writer, t reflect.Type, inline bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		writeQuery(w, t.Elem(), false)
	case reflect.Struct:
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			// nolint
			io.WriteString(w, "{")
		}
		for i := 0; i < t.NumField(); i++ {
			if i != 0 {
				// nolint
				io.WriteString(w, ",")
			}
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					// nolint
					io.WriteString(w, value)
				} else {
					// nolint
					//io.WriteString(w, ident.ParseMixedCaps(f.Name).ToLowerCamelCase())
					io.WriteString(w, f.Name)
				}
			}
			writeQuery(w, f.Type, inlineField)
		}
		if !inline {
			// nolint
			io.WriteString(w, "}")
		}
	}
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
