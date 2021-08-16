package graphqlparser

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"
)

const writeQueryDeepLimit = 10

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
	writeQuery(&buf, reflect.TypeOf(v), false, 0)
	return buf.String()
}

func writeQuery(w io.Writer, t reflect.Type, inline bool, depth uint) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		writeQuery(w, t.Elem(), false, depth)
	case reflect.Struct:

		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			// nolint
			io.WriteString(w, "{")
		}

		skipComma := false
		for i := 0; i < t.NumField(); i++ {
			hasStruct := false
			if depth > writeQueryDeepLimit {
				log.Println("You reached writeQueryDeepLimit constant")
				continue
			}
			if i != 0 && !skipComma {
				// nolint
				io.WriteString(w, ",")
			}
			skipComma = false
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					// nolint
					io.WriteString(w, value)
				} else {
					if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Ptr &&
						f.Type.Elem().Elem().Kind() == reflect.Struct ||
						f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Struct ||
						f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct && f.Type.Elem().Name() != "Time" ||
						f.Type.Kind() == reflect.Struct && f.Type.Name() != "Time" {
						hasStruct = true

						if depth+1 > writeQueryDeepLimit {
							log.Println("You reached writeQueryDeepLimit constant")
							skipComma = true
							continue
						}
					}
					// nolint
					//io.WriteString(w, ident.ParseMixedCaps(f.Name).ToLowerCamelCase())
					io.WriteString(w, f.Name)
				}
			}
			if hasStruct {
				writeQuery(w, f.Type, inlineField, depth+1)
			} else {
				writeQuery(w, f.Type, inlineField, depth)
			}
		}
		if !inline {
			// nolint
			io.WriteString(w, "}")
		}
	}
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
