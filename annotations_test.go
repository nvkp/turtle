package turtle_test

import (
	"testing"

	"github.com/nvkp/turtle"
	"github.com/nvkp/turtle/assert"
)

type tripleWithAnnotationValues struct {
	Subject    string `turtle:"subject"`
	Predicate  string `turtle:"predicate"`
	Object     string `turtle:"object"`
	Label      string `turtle:"label"`
	DataType   string `turtle:"datatype"`
	ObjectType string `turtle:"objecttype"`
}

type tripleWithAnnotationPointers struct {
	Subject    string  `turtle:"subject"`
	Predicate  string  `turtle:"predicate"`
	Object     string  `turtle:"object"`
	Label      *string `turtle:"label"`
	DataType   *string `turtle:"datatype"`
	ObjectType *string `turtle:"objecttype"`
}

var marshalWithAnnotationTestCases = map[string]struct {
	triples   interface{}
	expString string
	expErr    error
}{
	"one_triple_with_label_value": {
		triples: tripleWithAnnotationValues{
			Subject:    "http://example.org/person/Mark_Twain",
			Predicate:  "http://example.org/relation/name",
			Object:     "Huckleberry Finn",
			Label:      "en",
			ObjectType: "literal",
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"@en .
`,
	},
	"one_triple_with_datatype_value": {
		triples: tripleWithAnnotationValues{
			Subject:    "http://example.org/person/Mark_Twain",
			Predicate:  "http://example.org/relation/name",
			Object:     "Huckleberry Finn",
			DataType:   "xsd:string",
			ObjectType: "literal",
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"^^xsd:string .
`,
	},
	"one_triple_with_label_pointer": {
		triples: tripleWithAnnotationPointers{
			Subject:    "http://example.org/person/Mark_Twain",
			Predicate:  "http://example.org/relation/name",
			Object:     "Huckleberry Finn",
			Label:      ptr("en"),
			ObjectType: ptr("literal"),
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"@en .
`,
	},
	"one_triple_with_datatype_pointer": {
		triples: tripleWithAnnotationPointers{
			Subject:    "http://example.org/person/Mark_Twain",
			Predicate:  "http://example.org/relation/name",
			Object:     "Huckleberry Finn",
			DataType:   ptr("xsd:string"),
			ObjectType: ptr("literal"),
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"^^xsd:string .
`,
	},
	"slice_of_triples_with_annotations": {
		triples: []tripleWithAnnotationValues{
			{
				Subject:    "http://example.org/person/Mark_Twain",
				Predicate:  "http://example.org/relation/name",
				Object:     "Huckleberry Finn",
				DataType:   "xsd:string",
				ObjectType: "literal",
			},
			{
				Subject:    "http://example.org/person/Mark_Twain",
				Predicate:  "http://example.org/relation/name",
				Object:     "Huckleberry Finn",
				Label:      "en",
				ObjectType: "literal",
			},
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"^^xsd:string, "Huckleberry Finn"@en .
`,
	},
}

func TestMarshalWithAnnotations(t *testing.T) {
	for name, tc := range marshalWithAnnotationTestCases {
		t.Run(name, func(t *testing.T) {
			b, err := turtle.Marshal(tc.triples)
			assert.Equal(t, tc.expString, string(b), "Marshal function should have returned a correct byte data")
			assert.ErrorIs(t, err, tc.expErr, "Marshal function should have returned a correct error")
		})
	}
}

func TestUnmarshalStructWithLabel(t *testing.T) {
	var target tripleWithAnnotationValues
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"@en .`)
	expected := tripleWithAnnotationValues{
		Subject:    "http://example.org/person/Mark_Twain",
		Predicate:  "http://example.org/relation/name",
		Object:     "Huckleberry Finn",
		Label:      "en",
		ObjectType: "literal",
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target triple")
}

func TestUnmarshalStructWithDataType(t *testing.T) {
	var target tripleWithAnnotationValues
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"^^xsd:string .`)
	expected := tripleWithAnnotationValues{
		Subject:    "http://example.org/person/Mark_Twain",
		Predicate:  "http://example.org/relation/name",
		Object:     "Huckleberry Finn",
		DataType:   "xsd:string",
		ObjectType: "literal",
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target triple")
}

func TestUnmarshalStructWithLabelPointer(t *testing.T) {
	var target tripleWithAnnotationPointers
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/name> "Huckleberry Finn"@en .`)
	expected := tripleWithAnnotationPointers{
		Subject:    "http://example.org/person/Mark_Twain",
		Predicate:  "http://example.org/relation/name",
		Object:     "Huckleberry Finn",
		Label:      ptr("en"),
		ObjectType: ptr("literal"),
	}

	err := turtle.Unmarshal(data, &target)

	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target triple")
}
