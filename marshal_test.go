package turtle_test

import (
	"strings"
	"testing"

	"github.com/nvkp/turtle"
	"github.com/nvkp/turtle/assert"
)

type triple struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}

type tripleWithPointers struct {
	Subject   *string `turtle:"subject"`
	Predicate *string `turtle:"predicate"`
	Object    *string `turtle:"object"`
}

type subject string
type predicate string
type object string

type namedTypeTriple struct {
	s subject   `turtle:"subject"`
	p predicate `turtle:"predicate"`
	o object    `turtle:"object"`
}

type tripleWithMetadata struct {
	Subject   string            `turtle:"subject"`
	Predicate string            `turtle:"predicate"`
	Object    string            `turtle:"object"`
	Prefixes  map[string]string `turtle:"prefix"`
	Base      string            `turtle:"base"`
}

var marshalTestCases = map[string]struct {
	triples   interface{}
	expString string
	expErr    error
}{
	"one_triple": {
		triples: triple{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author",
			Object:    "http://example.org/books/Huckleberry_Finn",
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
`,
	},
	"one_triple_pointer": {
		triples: &triple{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author",
			Object:    "http://example.org/books/Huckleberry_Finn",
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
`,
	},
	"named_type_triple": {
		triples: namedTypeTriple{
			s: "http://example.org/person/Mark_Twain",
			p: "http://example.org/relation/author",
			o: "http://example.org/books/Huckleberry_Finn",
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
`,
	},
	"slice": {
		triples: []triple{
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://www.perceive.net/schemas/relationship/enemyOf",
				Object:    "http://example.org/spiderman",
			},
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				Object:    "http://xmlns.com/foaf/0.1/Person",
			},
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Green Goblin",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://www.perceive.net/schemas/relationship/enemyOf",
				Object:    "http://example.org/green-goblin",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				Object:    "http://xmlns.com/foaf/0.1/Person",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Spiderman",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Человек-паук",
			},
		},
		expString: `<http://example.org/green-goblin> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Green Goblin" .
<http://example.org/spiderman> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Spiderman", "Человек-паук" .
`,
	},
	"slice_pointer": {
		triples: &[]triple{
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://www.perceive.net/schemas/relationship/enemyOf",
				Object:    "http://example.org/spiderman",
			},
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				Object:    "http://xmlns.com/foaf/0.1/Person",
			},
			{
				Subject:   "http://example.org/green-goblin",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Green Goblin",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://www.perceive.net/schemas/relationship/enemyOf",
				Object:    "http://example.org/green-goblin",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				Object:    "http://xmlns.com/foaf/0.1/Person",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Spiderman",
			},
			{
				Subject:   "http://example.org/spiderman",
				Predicate: "http://xmlns.com/foaf/0.1/name",
				Object:    "Человек-паук",
			},
		},
		expString: `<http://example.org/green-goblin> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Green Goblin" .
<http://example.org/spiderman> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Spiderman", "Человек-паук" .
`,
	},
	"named_type_slice": {
		triples: []namedTypeTriple{
			{
				s: "http://example.org/green-goblin",
				p: "http://www.perceive.net/schemas/relationship/enemyOf",
				o: "http://example.org/spiderman",
			},
			{
				s: "http://example.org/green-goblin",
				p: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				o: "http://xmlns.com/foaf/0.1/Person",
			},
			{
				s: "http://example.org/green-goblin",
				p: "http://xmlns.com/foaf/0.1/name",
				o: "Green Goblin",
			},
			{
				s: "http://example.org/spiderman",
				p: "http://www.perceive.net/schemas/relationship/enemyOf",
				o: "http://example.org/green-goblin",
			},
			{
				s: "http://example.org/spiderman",
				p: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
				o: "http://xmlns.com/foaf/0.1/Person",
			},
			{
				s: "http://example.org/spiderman",
				p: "http://xmlns.com/foaf/0.1/name",
				o: "Spiderman",
			},
			{
				s: "http://example.org/spiderman",
				p: "http://xmlns.com/foaf/0.1/name",
				o: "Человек-паук",
			},
		},
		expString: `<http://example.org/green-goblin> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Green Goblin" .
<http://example.org/spiderman> 
	<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
	<http://xmlns.com/foaf/0.1/name> "Spiderman", "Человек-паук" .
`,
	},
	"triple_struct_with_pointers": {
		triples: &tripleWithPointers{
			Subject:   ptr("http://example.org/person/Mark_Twain"),
			Predicate: ptr("http://example.org/relation/author"),
			Object:    ptr("http://example.org/books/Huckleberry_Finn"),
		},
		expString: `<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
`,
	},
	"nil_input": {
		triples:   nil,
		expString: ``,
		expErr:    turtle.ErrInvalidValueType,
	},
	"no_subject_specified": {
		triples: triple{
			Predicate: "http://example.org/relation/author",
			Object:    "http://example.org/books/Huckleberry_Finn",
		},
		expString: ``,
		expErr:    turtle.ErrNoSubjectSpecified,
	},
	"no_predicate_specified": {
		triples: triple{
			Subject: "http://example.org/person/Mark_Twain",
			Object:  "http://example.org/books/Huckleberry_Finn",
		},
		expString: ``,
		expErr:    turtle.ErrNoPredicateSpecified,
	},
	"no_object_specified": {
		triples: triple{
			Subject:   "http://example.org/person/Mark_Twain",
			Predicate: "http://example.org/relation/author",
		},
		expString: ``,
		expErr:    turtle.ErrNoObjectSpecified,
	},
}

func TestMarshal(t *testing.T) {
	for name, tc := range marshalTestCases {
		t.Run(name, func(t *testing.T) {
			b, err := turtle.Marshal(tc.triples)
			assert.Equal(t, tc.expString, string(b), "Marshal function should have returned a correct byte data")
			assert.ErrorIs(t, err, tc.expErr, "Marshal function should have returned a correct error")
		})
	}
}

func TestMarshalOptions(t *testing.T) {
	trip := triple{
		Subject:   "http://example.org/person/Mark_Twain",
		Predicate: "http://example.org/relation/author",
		Object:    "http://example.org/books/Huckleberry_Finn",
	}

	c := turtle.Config{
		Base: "http://example.org",
		Prefixes: map[string]string{
			"book": "http://example.org/books/",
		},
	}

	out, err := c.Marshal(trip)
	assert.NoError(t, err, "no error was expected")
	assert.Equal(t, strings.TrimSpace(string(out)), strings.TrimSpace(`
@base <http://example.org> .
@prefix book: <http://example.org/books/> .
<http://example.org/person/Mark_Twain> <http://example.org/relation/author> <http://example.org/books/Huckleberry_Finn> .
`), "output was not equal")

	c.ResolveURLs = true

	out, err = c.Marshal(trip)
	assert.NoError(t, err, "no error was expected")
	assert.Equal(t, strings.TrimSpace(string(out)), strings.TrimSpace(`
@base <http://example.org> .
@prefix book: <http://example.org/books/> .
</person/Mark_Twain> </relation/author> book:Huckleberry_Finn .
`), "output was not equal")

	// check for weird rdf-isms like using a blank anchor as a prefix

	c = turtle.Config{
		Base: "http://example.org",
		Prefixes: map[string]string{
			"book": "http://example.org/books#",
		},
		ResolveURLs: true,
	}

	trip = triple{
		Subject:   "http://example.org/person/Mark_Twain",
		Predicate: "http://example.org/relation/author",
		Object:    "http://example.org/books#Huckleberry_Finn",
	}

	out, err = c.Marshal(trip)
	assert.NoError(t, err, "no error was expected")
	assert.Equal(t, strings.TrimSpace(string(out)), strings.TrimSpace(`
@base <http://example.org> .
@prefix book: <http://example.org/books#> .
</person/Mark_Twain> </relation/author> book:Huckleberry_Finn .
`), "output was not equal")
}

func TestMarshalPrefixBase(t *testing.T) {
	_, err := turtle.Marshal(tripleWithMetadata{
		Subject:   "http://example.org/Mark_Twain",
		Predicate: "a",
		Object:    "foaf:Person",
		Prefixes: map[string]string{
			"foaf": "http://xmlns.com/foaf/0.1",
		},
		Base: "http://example.org",
	})

	assert.NoError(t, err, "no error expected marshaling with base and prefixes")
}
