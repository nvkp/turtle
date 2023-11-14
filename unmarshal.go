package turtle

import (
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	s := newScanner(data)
	return decode(s, reflect.ValueOf(v)) // TODO panics when not pointer, what does goccy do?
}

/*
if typ == nil || typ.Kind() != reflect.Ptr || p == 0 {
		return &InvalidUnmarshalError{Type: runtime.RType2Type(typ)}
	}
*/

func decode(s *scanner, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		return decode(s, v.Elem())
	case reflect.Slice:
		for s.next() {
			item := reflect.New(v.Type().Elem()).Elem()
			// TODO check item is struct
			decodeStruct(s, item)
			v.Set(reflect.Append(v, item))
		}
	case reflect.Struct:
		ok := s.next()
		if !ok {
			return nil
		}
		decodeStruct(s, v)
	}

	return nil
}

func decodeStruct(s *scanner, v reflect.Value) {
	t := s.triple()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !(field.Kind() == reflect.String) {
			continue
		}

		tag := v.Type().Field(i).Tag.Get("turtle")
		if tag == "" {
			continue
		}

		if tag == "subject" {
			field.SetString(t[0])
		}

		if tag == "predicate" {
			field.SetString(t[1])
		}

		if tag == "object" {
			field.SetString(t[2])
		}
	}
}
