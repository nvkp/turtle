package turtle

import (
	"fmt"
	"reflect"

	"errors"

	"github.com/nvkp/turtle/graph"
)

var (
	ErrInvalidValueType     = errors.New("invalid value's type")
	ErrNoPointerValue       = errors.New("value not a pointer")
	ErrNilValue             = errors.New("value is nil")
	ErrNoSubjectSpecified   = errors.New("no subject tag specified in struct")
	ErrNoPredicateSpecified = errors.New("no predicate tag specified in struct")
	ErrNoObjectSpecified    = errors.New("no object tag specified in struct")
)

// Marshal TODO comment
func Marshal(v interface{}) ([]byte, error) {
	g := make(graph.Graph)
	if err := marshal(g, reflect.ValueOf(v)); err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return g.Bytes()
}

func marshal(g graph.Graph, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		// if value is pointer marhal the pointed value
		return marshal(g, v.Elem())
	case reflect.Array, reflect.Slice:
		// if value is iterable iterate over value and marshal each element
		for i := 0; i < v.Len(); i++ {
			if err := marshal(g, v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		// if value is struct go look into the struct's fields
		return marshalStruct(g, v)
	default:
		return ErrInvalidValueType
	}

	return nil
}

func marshalStruct(g graph.Graph, v reflect.Value) error {
	var t [3]string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// get the tag value
		tag := v.Type().Field(i).Tag.Get("turtle")
		if tag == "" {
			continue
		}

		// pick correct part based on the tag value
		var part int
		switch tag {
		case "subject":
			part = subject
		case "predicate":
			part = predicate
		case "object":
			part = object
		}

		var word string
		// if field is string use its value
		if field.Kind() == reflect.String {
			word = field.String()
		}
		// is field is pointer to string use the pointed value
		if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String {
			word = field.Elem().String()
		}

		// fill the word to correct part of the triple
		t[part] = word
	}

	if t[subject] == "" {
		return ErrNoSubjectSpecified
	}

	if t[predicate] == "" {
		return ErrNoPredicateSpecified
	}

	if t[object] == "" {
		return ErrNoObjectSpecified
	}

	// accept the extracted triple to graph
	g.Accept(t)

	return nil
}
