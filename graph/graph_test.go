package graph_test

import (
	"testing"

	"github.com/nvkp/turtle/assert"
	"github.com/nvkp/turtle/graph"
)

var graphTestCases = map[string]struct {
	triples  [][3]string
	expected []byte
}{
	"simple_graph": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"c", "d", "e"},
		},
		expected: []byte(`<a> <b> <c> .
<c> <d> <e> .
`),
	},
	"subject_with_two_predicates": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "c", "e"},
		},
		expected: []byte(`<a> 
	<b> <c> ;
	<c> <e> .
`),
	},
	"predicate_with_two_objects": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "b", "d"},
		},
		expected: []byte(`<a> <b> <c>, <d> .
`),
	},
	"two_predicates_with_two_objects": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "b", "d"},
			{"a", "e", "c"},
			{"a", "e", "d"},
		},
		expected: []byte(`<a> 
	<b> <c>, <d> ;
	<e> <c>, <d> .
`),
	},
}

func TestGraph(t *testing.T) {
	for name, tc := range graphTestCases {
		t.Run(name, func(t *testing.T) {
			g := graph.New()

			for _, triple := range tc.triples {
				_ = g.Accept(triple)
			}

			b, _ := g.Bytes()

			assert.Equal(t, tc.expected, b, "accept method should have created a correct graph structure")
		})
	}
}
