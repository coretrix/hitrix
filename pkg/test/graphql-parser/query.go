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
		_, err := io.WriteString(&buf, "$")
		if err != nil {
			panic(err)
		}

		_, err = io.WriteString(&buf, k)
		if err != nil {
			panic(err)
		}

		_, err = io.WriteString(&buf, ":")
		if err != nil {
			panic(err)
		}

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
		_, err := io.WriteString(w, "[")
		if err != nil {
			panic(err)
		}

		writeArgumentType(w, t.Elem(), true)

		_, err = io.WriteString(w, "]")
		if err != nil {
			panic(err)
		}
	default:
		_, err := io.WriteString(w, strings.Title(t.Name()))
		if err != nil {
			panic(err)
		}
	}

	if value {
		_, err := io.WriteString(w, "!")
		if err != nil {
			panic(err)
		}
	}
}

func query(v interface{}) string {
	var buf bytes.Buffer
	shownError := false
	writeQuery(&buf, reflect.TypeOf(v), false, 0, &shownError)

	return buf.String()
}

func writeQuery(w io.Writer, t reflect.Type, inline bool, depth uint, shownError *bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		writeQuery(w, t.Elem(), false, depth, shownError)
	case reflect.Struct:
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			_, err := io.WriteString(w, "{")
			if err != nil {
				panic(err)
			}
		}

		skipComma := false
		for i := 0; i < t.NumField(); i++ {
			hasStruct := false
			if depth > writeQueryDeepLimit {
				if !*shownError {
					*shownError = true
					log.Println("You reached writeQueryDeepLimit constant")
				}

				continue
			}
			if i != 0 && !skipComma {
				_, err := io.WriteString(w, ",")
				if err != nil {
					panic(err)
				}
			}
			skipComma = false
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			fieldName, hasJSON := f.Tag.Lookup("json")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					_, err := io.WriteString(w, value)
					if err != nil {
						panic(err)
					}
				} else {
					if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Ptr &&
						f.Type.Elem().Elem().Kind() == reflect.Struct ||
						f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Struct ||
						f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct && f.Type.Elem().Name() != "Time" ||
						f.Type.Kind() == reflect.Struct && f.Type.Name() != "Time" {
						hasStruct = true

						if depth+1 > writeQueryDeepLimit {
							if !*shownError {
								*shownError = true
								log.Println("You reached writeQueryDeepLimit constant")
							}
							skipComma = true

							continue
						}
					}

					//io.WriteString(w, ident.ParseMixedCaps(f.Name).ToLowerCamelCase())
					if hasJSON && fieldName != "-" {
						_, err := io.WriteString(w, fieldName)
						if err != nil {
							panic(err)
						}
					} else {
						_, err := io.WriteString(w, f.Name)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			if hasStruct {
				writeQuery(w, f.Type, inlineField, depth+1, shownError)
			} else {
				writeQuery(w, f.Type, inlineField, depth, shownError)
			}
		}
		if !inline {
			_, err := io.WriteString(w, "}")
			if err != nil {
				panic(err)
			}
		}
	}
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
