package turtle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var graphTestCases = map[string]struct {
	triples  [][3]string
	expected graph
}{
	"simple_graph": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"c", "d", "e"},
		},
		expected: graph{
			"a": {
				"b": {"c"},
			},
			"c": {
				"d": {"e"},
			},
		},
	},
	"subject_with_two_predicates": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "c", "e"},
		},
		expected: graph{
			"a": {
				"b": {"c"},
				"c": {"e"},
			},
		},
	},
	"predicate_with_two_objects": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "b", "d"},
		},
		expected: graph{
			"a": {
				"b": {"c", "d"},
			},
		},
	},
	"two_predicates_with_two_objects": {
		triples: [][3]string{
			{"a", "b", "c"},
			{"a", "b", "d"},
			{"a", "e", "c"},
			{"a", "e", "d"},
		},
		expected: graph{
			"a": {
				"b": {"c", "d"},
				"e": {"c", "d"},
			},
		},
	},
}

func TestGraph(t *testing.T) {
	for name, tc := range graphTestCases {
		t.Run(name, func(t *testing.T) {
			g := make(graph)

			for _, triple := range tc.triples {
				_ = g.accept(triple)
			}

			assert.Equal(t, tc.expected, g, "accept method should have created a correct graph structure")
		})
	}
}
