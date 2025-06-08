package turtle

import (
	"reflect"

	"errors"

	"github.com/nvkp/turtle/graph"
)

var (
	// ErrInvalidValueType is returned by Marshal when an invalid value is passed
	ErrInvalidValueType = errors.New("invalid value's type")
	// ErrNoSubjectSpecified is returned by Marshal when the provided struct does not contain a field with a turtle:"subject" tag
	ErrNoSubjectSpecified = errors.New("no subject tag specified in struct")
	// ErrNoPredicateSpecified is returned by Marshal when the provided struct does not contain a field with a turtle:"predicate" tag
	ErrNoPredicateSpecified = errors.New("no predicate tag specified in struct")
	// ErrNoObjectSpecified is returned by Marshal when the provided struct does not contain a field with a turtle:"object" tag
	ErrNoObjectSpecified = errors.New("no object tag specified in struct")
)

// Marshal serializes the provided data structure into RDF Turtle format.
// The function accepts the to-be-serialized data as an empty interface
// and returns the byte slice with the result and possible error value.
// It is able to handle single struct, struct, a slice, an array or
// a pointer to all three.
//
// The fields of the structs passed to the function have to be annotated
// by Golang tag `turtle` defining which of the fields correspond to
// which part of the RDF triple (either "subject", "predicate" or "object").
//
// The compact version of the Turtle format is used. The resulting Turtle
// triples are sorted alphabetically first by subjects, then by predicates
// and then by objects.
func Marshal(v interface{}) ([]byte, error) {
	return (&Config{}).Marshal(v)
}

func marshal(g *graph.Graph, v reflect.Value) error {
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

func marshalStruct(g *graph.Graph, v reflect.Value) error {
	var t [6]string

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
		case "label":
			part = label
		case "datatype":
			part = datatype
		case "objecttype":
			part = objecttype
		case "base", "prefix":
			continue
		}

		var word string
		// if field is string use its value
		if field.Kind() == reflect.String {
			word = field.String()
		}
		// is field is pointer to string use the pointed value
		if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String && field.Elem().Kind() == reflect.String {
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
	g.AcceptWithAnnotations(t)

	return nil
}
