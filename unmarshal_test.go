package turtle_test

import (
	"testing"

	"github.com/nvkp/turtle"
	"github.com/nvkp/turtle/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestUnmarshalStruct(t *testing.T) {
	var target triple
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .`)
	expected := triple{
		Subject:   "http://example.org/person/Mark_Twain",
		Predicate: "http://example.org/relation/author",
		Object:    "http://example.org/books/Huckleberry_Finn",
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target triple")
}

func TestUnmarshalSlice(t *testing.T) {
	target := make([]triple, 0)
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
	<http://example.org/person/Mark_Twain> <http://example.org/relation/author2> <http://example.org/books/Huckleberry_Finn2> .`)
	expected := []triple{
		{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author",
			Object:    "http://example.org/books/Huckleberry_Finn",
		},
		{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author2",
			Object:    "http://example.org/books/Huckleberry_Finn2",
		},
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target slice")
}

func TestUnmarshalCompact(t *testing.T) {
	target := make([]triple, 0)
	data := []byte(`<http://example.org/green-goblin>
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Green Goblin".<http://example.org/spiderman>
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin>;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person>;
	<http://xmlns.com/foaf/0.1/name> "Spiderman", "Человек-паук" .`)
	expected := []triple{
		{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
		{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
		{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		{"http://example.org/spiderman", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/green-goblin"},
		{"http://example.org/spiderman", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
		{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Spiderman"},
		{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Человек-паук"},
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target slice")
}

func TestUnmarshalSliceOfPointers(t *testing.T) {
	target := make([]*triple, 0)
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
	<http://example.org/person/Mark_Twain> <http://example.org/relation/author2> <http://example.org/books/Huckleberry_Finn2> .`)
	expected := []*triple{
		{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author",
			Object:    "http://example.org/books/Huckleberry_Finn",
		},
		{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author2",
			Object:    "http://example.org/books/Huckleberry_Finn2",
		},
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target slice")
}

func TestUnmarshalSliceStructsWithPointers(t *testing.T) {
	target := make([]tripleWithPointers, 0)
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
	<http://example.org/person/Mark_Twain> <http://example.org/relation/author2> <http://example.org/books/Huckleberry_Finn2> .`)
	expected := []tripleWithPointers{
		{
			Subject:   ptr("http://example.org/person/Mark_Twain"),
			Predicate: ptr("http://example.org/relation/author"),
			Object:    ptr("http://example.org/books/Huckleberry_Finn"),
		},
		{
			Subject:   ptr("http://example.org/person/Mark_Twain"),
			Predicate: ptr("http://example.org/relation/author2"),
			Object:    ptr("http://example.org/books/Huckleberry_Finn2"),
		},
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "function Unmarshal should have returned no error")
	assert.Equal(t, expected, target, "function Unmarshal should have assigned correct values to the target slice")
}

func TestUnmarshalNil(t *testing.T) {
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .`)

	err := turtle.Unmarshal(data, nil)
	assert.ErrorIs(t, err, turtle.ErrNilValue, "function Unmarshal should have returned correct error")
}

func TestUnmarshalNotAPointer(t *testing.T) {
	var target triple
	data := []byte(`<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .`)

	err := turtle.Unmarshal(data, target)
	assert.ErrorIs(t, err, turtle.ErrNoPointerValue, "function Unmarshal should have returned correct error")
}

func TestUnmarshalBaseAndPrefixes(t *testing.T) {
	target := make([]tripleWithMetadata, 0)
	data := []byte(`
@base <http://example.org/> .
@prefix books: <https://amazon.com/> .
</person/Mark_Twain> </relation/author> <books:Huckleberry_Finn> .`)

	expected := []tripleWithMetadata{
		{
			Base:      "http://example.org/",
			Prefixes:  map[string]string{"books": "https://amazon.com/"},
			Subject:   "/person/Mark_Twain",
			Predicate: "/relation/author",
			Object:    "books:Huckleberry_Finn",
		},
	}

	err := turtle.Unmarshal(data, &target)
	assert.NoError(t, err, "got an error unmarshaling turtle with base and prefixes")
	assert.Equal(t, expected, target, "not equal to expected data")
}
