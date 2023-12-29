package turtle

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nvkp/turtle/scanner"
)

const (
	subject int = iota
	predicate
	object
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
	if v == nil {
		return ErrNilValue
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return ErrNoPointerValue
	}

	err := unmarshal(scanner.New(data), rv)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
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
		return unmarshalStruct(s, v)
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

		switch itemType.Kind() {
		case reflect.Pointer:
			// if slice contains pointers, create item as
			// reflect pointer to element zero value and call unmarshalStruct on
			// the pointer's element
			item = reflect.New(itemType.Elem())
			err = unmarshalStruct(s, item.Elem())
		case reflect.Struct:
			// if slice contains structs, create item as
			// zero value struct and call unmarshalStruct on it
			item = reflect.New(itemType).Elem()
			err = unmarshalStruct(s, item)
		default:
			return errors.New("invalid slice's item type")
		}

		if err != nil {
			return err
		}

		if !v.CanSet() {
			return errors.New("cannot append to slice")
		}

		// append the filled item (either struct or pointer to struct) to the slice
		v.Set(reflect.Append(v, item))
	}

	return nil
}

func unmarshalStruct(s *scanner.Scanner, v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return errors.New("value not struct")
	}
	t := s.Triple()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			return errors.New("field cannot be changed")
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
		}
		word = t[part]

		// if field is string set value
		if field.Kind() == reflect.String {
			field.SetString(word)
		}

		// if field is pointer
		if field.Kind() == reflect.Pointer {
			pointerType := field.Type().Elem()
			// check that the field is pointer to string
			if pointerType.Kind() != reflect.String {
				continue
			}

			// create new reflect value of pointer to string
			stringValue := reflect.New(pointerType)
			// set value to the pointed string
			stringValue.Elem().SetString(word)
			// set the pointer as value of the struct field
			field.Set(stringValue)
		}
	}

	return nil
}
