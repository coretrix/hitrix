package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func UnmarshalGraphQL(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := (&decoder{tokenizer: dec}).Decode(v)
	if err != nil {
		return err
	}
	tok, err := dec.Token()
	switch err {
	case io.EOF:

		return nil
	case nil:
		return fmt.Errorf("invalid token '%v' after top-level value", tok)
	default:
		return err
	}
}

type decoder struct {
	vs         [][]reflect.Value
	parseState []json.Delim
	tokenizer  interface {
		Token() (json.Token, error)
	}
}

func (d *decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot decode into non-pointer %T", v)
	}
	d.vs = [][]reflect.Value{{rv.Elem()}}
	return d.decode()
}

func (d *decoder) decode() error {
	for len(d.vs) > 0 {
		tok, err := d.tokenizer.Token()
		if err == io.EOF {
			return errors.New("unexpected end of JSON input")
		} else if err != nil {
			return err
		}

		switch {
		case d.state() == '{' && tok != json.Delim('}'):
			key, ok := tok.(string)
			if !ok {
				return errors.New("unexpected non-key in JSON input")
			}
			someFieldExist := false
			for i := range d.vs {
				v := d.vs[i][len(d.vs[i])-1]
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				var f reflect.Value
				if v.Kind() == reflect.Struct {
					f = fieldByGraphQLName(v, key)
					if f.IsValid() {
						someFieldExist = true
					}
				}
				d.vs[i] = append(d.vs[i], f)
			}
			if !someFieldExist {
				return fmt.Errorf("struct field for %q doesn't exist in any of %v places to unmarshal", key, len(d.vs))
			}

			tok, err = d.tokenizer.Token()
			if err == io.EOF {
				return errors.New("unexpected end of JSON input")
			} else if err != nil {
				return err
			}

		case d.state() == '[' && tok != json.Delim(']'):
			someSliceExist := false
			for i := range d.vs {
				v := d.vs[i][len(d.vs[i])-1]
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				var f reflect.Value
				if v.Kind() == reflect.Slice {
					v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
					f = v.Index(v.Len() - 1)
					someSliceExist = true
				}
				d.vs[i] = append(d.vs[i], f)
			}
			if !someSliceExist {
				return fmt.Errorf("slice doesn't exist in any of %v places to unmarshal", len(d.vs))
			}
		}

		switch tok := tok.(type) {
		case string, json.Number, bool, nil:
			for i := range d.vs {
				v := d.vs[i][len(d.vs[i])-1]
				if !v.IsValid() {
					continue
				}
				err := unmarshalValue(tok, v)
				if err != nil {
					return err
				}
			}
			d.popAllVs()

		case json.Delim:
			switch tok {
			case '{':

				d.pushState(tok)

				frontier := make([]reflect.Value, len(d.vs))
				for i := range d.vs {
					v := d.vs[i][len(d.vs[i])-1]
					frontier[i] = v
					if v.Kind() == reflect.Ptr && v.IsNil() {
						v.Set(reflect.New(v.Type().Elem()))
					}
				}

				for len(frontier) > 0 {
					v := frontier[0]
					frontier = frontier[1:]
					if v.Kind() == reflect.Ptr {
						v = v.Elem()
					}
					if v.Kind() != reflect.Struct {
						continue
					}
					for i := 0; i < v.NumField(); i++ {
						if isGraphQLFragment(v.Type().Field(i)) || v.Type().Field(i).Anonymous {
							d.vs = append(d.vs, []reflect.Value{v.Field(i)})
							frontier = append(frontier, v.Field(i))
						}
					}
				}
			case '[':

				d.pushState(tok)

				for i := range d.vs {
					v := d.vs[i][len(d.vs[i])-1]
					if v.Kind() == reflect.Ptr {
						v = v.Elem()
					}
					if v.Kind() != reflect.Slice {
						continue
					}
					v.Set(reflect.MakeSlice(v.Type(), 0, 0))
				}
			case '}', ']':
				d.popAllVs()
				d.popState()
			default:
				return errors.New("unexpected delimiter in JSON input")
			}
		default:
			return errors.New("unexpected token in JSON input")
		}
	}
	return nil
}

func (d *decoder) pushState(s json.Delim) {
	d.parseState = append(d.parseState, s)
}

func (d *decoder) popState() {
	d.parseState = d.parseState[:len(d.parseState)-1]
}

func (d *decoder) state() json.Delim {
	if len(d.parseState) == 0 {
		return 0
	}
	return d.parseState[len(d.parseState)-1]
}

func (d *decoder) popAllVs() {
	var nonEmpty [][]reflect.Value
	for i := range d.vs {
		d.vs[i] = d.vs[i][:len(d.vs[i])-1]
		if len(d.vs[i]) > 0 {
			nonEmpty = append(nonEmpty, d.vs[i])
		}
	}
	d.vs = nonEmpty
}

func fieldByGraphQLName(v reflect.Value, name string) reflect.Value {
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).PkgPath != "" {
			continue
		}
		if hasGraphQLName(v.Type().Field(i), name) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

func hasGraphQLName(f reflect.StructField, name string) bool {
	value, ok := f.Tag.Lookup("graphql")
	if !ok {
		return strings.EqualFold(f.Name, name)
	}
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "...") {
		return false
	}
	if i := strings.Index(value, "("); i != -1 {
		value = value[:i]
	}
	if i := strings.Index(value, ":"); i != -1 {
		value = value[:i]
	}
	return strings.TrimSpace(value) == name
}

func isGraphQLFragment(f reflect.StructField) bool {
	value, ok := f.Tag.Lookup("graphql")
	if !ok {
		return false
	}
	value = strings.TrimSpace(value)
	return strings.HasPrefix(value, "...")
}

func unmarshalValue(value json.Token, v reflect.Value) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v.Addr().Interface())
}
