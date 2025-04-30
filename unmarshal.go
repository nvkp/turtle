package turtle

import (
	"errors"
	"reflect"

	"github.com/nvkp/turtle/scanner"
)

const (
	subject int = iota
	predicate
	object
	label
	datatype
)

var (
	// ErrNoPointerValue is returned by Unmarshal function when the passed value is not a pointer
	ErrNoPointerValue = errors.New("value not a pointer")
	// ErrNilValue is returned by Unmarshal function when the passed value is nil
	ErrNilValue = errors.New("value is nil")
)

// Unmarshal parses Turtle data. It accepts a byte slice of
// the turtle data and also a target as a pointer to a
// struct or to a slice/array of structs that have fields
// annotated by tags turtle defining which field of the
// struct corresponds to which part of the RDF triple.
//
// The function accepts the compact version of Turtle
// just as the N-triples version of the format where each
// row corresponds to a single triple. It reads @base
// and @prefix forms and extends the IRIs that are filled
// in the target structure with them. It ignores Turtle
// comments, labels and data types. The keyword a gets
// replaced by http://www.w3.org/1999/02/22-rdf-syntax-ns#type IRI.
func Unmarshal(data []byte, v interface{}) error {
	return (&Config{}).Unmarshal(data, v)
}

func unmarshal(s *scanner.Scanner, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		return unmarshal(s, v.Elem())
	case reflect.Slice:
		return unmarshalSlice(s, v)
	case reflect.Struct:
		ok := s.Next()
		if !ok {
			return nil
		}
		err, _ := unmarshalStruct(s, v)
		return err
	}

	return nil
}

func unmarshalSlice(s *scanner.Scanner, v reflect.Value) error {
	if v.Kind() != reflect.Slice {
		return errors.New("value not a slice")
	}
	// get type of the elements of the slice
	itemType := v.Type().Elem()

	for s.Next() {
		var item reflect.Value
		var err error
		var ok bool

		switch itemType.Kind() {
		case reflect.Pointer:
			// if slice contains pointers, create item as
			// reflect pointer to element zero value and call unmarshalStruct on
			// the pointer's element
			item = reflect.New(itemType.Elem())
			err, ok = unmarshalStruct(s, item.Elem())
		case reflect.Struct:
			// if slice contains structs, create item as
			// zero value struct and call unmarshalStruct on it
			item = reflect.New(itemType).Elem()
			err, ok = unmarshalStruct(s, item)
		default:
			return errors.New("invalid slice's item type")
		}

		if err != nil {
			return err
		}

		if !ok {
			break
		}

		if !v.CanSet() {
			return errors.New("cannot append to slice")
		}

		// append the filled item (either struct or pointer to struct) to the slice
		v.Set(reflect.Append(v, item))
	}

	return nil
}

func unmarshalStruct(s *scanner.Scanner, v reflect.Value) (error, bool) {
	if v.Kind() != reflect.Struct {
		return errors.New("value not struct"), false
	}

	var t [5]string

	// prevent empty structs from being generated from pragmas
outer:
	for {
		t = s.TripleWithAnnotations()
		for _, k := range t {
			if k != "" {
				break outer
			}
		}
		if !s.Next() {
			return nil, false
		}
	}

	numField := v.NumField()
	_ = numField
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			return errors.New("field cannot be changed"), false
		}

		// get the tag value
		tag := v.Type().Field(i).Tag.Get("turtle")
		if tag == "" {
			continue
		}

		// pick correct value from current triple based on the tag value
		var word string
		var part int
		switch tag {
		case "subject":
			part = subject
		case "predicate":
			part = predicate
		case "object":
			part = object
		case "label":
			part = label
		case "datatype":
			part = datatype
		case "base", "prefix":
			part = -1
		}

		if part >= 0 {
			word = t[part]
		}

		// if field is string set value
		if field.Kind() == reflect.String && tag != "prefix" {
			if tag == "base" {
				field.SetString(s.Base())
			} else {
				field.SetString(word)
			}
		} else if field.Kind() == reflect.Map && tag == "prefix" && isMap(field.Type()) {
			field.Set(reflect.ValueOf(s.Prefixes()))
		} else if field.Kind() == reflect.Pointer {
			pointerType := field.Type().Elem()
			// create new reflect value of pointer to string
			value := reflect.New(pointerType)

			// check that the field is pointer to string
			if pointerType.Kind() == reflect.String {

				// omit empty strings
				if len(word) == 0 {
					continue
				}

				// set value to the pointed string
				value.Elem().SetString(word)
				// set the pointer as value of the struct field
				field.Set(value)
			} else if isMap(pointerType) {
				field.Set(value)
			}
		}
	}

	return nil, true
}

func isMap(value reflect.Type) bool {
	return value.Key().Kind() == reflect.String && value.Elem().Kind() == reflect.String
}
