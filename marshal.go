package turtle

import (
	"errors"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	g := make(graph)
	if err := encode(g, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return g.bytes()
}

func encode(g graph, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		return encode(g, v.Elem())
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := encode(g, v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		var s, p, o string

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !(v.Field(i).Kind() == reflect.String) {
				continue
			}

			tag := field.Tag.Get("turtle")
			if tag == "" {
				continue
			}

			if tag == "subject" {
				s = v.Field(i).String()
			}

			if tag == "predicate" {
				p = v.Field(i).String()
			}

			if tag == "object" {
				o = v.Field(i).String()
			}
		}

		if s == "" || p == "" || o == "" {
			return nil // TODO ?
		}

		// accept to graph
		g.accept([3]string{s, p, o})
	default:
		return errors.New("invalid") // TODO organize errors
	}

	return nil
}
