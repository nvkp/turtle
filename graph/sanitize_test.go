package graph

import (
	"testing"

	"github.com/nvkp/turtle/assert"
)

var sanitizesTestCases = map[string]struct {
	str      string
	expected string
}{
	"empty_string": {
		str:      "",
		expected: "",
	},
	"iri": {
		str:      "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
		expected: "<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
	},
	"blank_node": {
		str:      "_:b23",
		expected: "_:b23",
	},
	"literal": {
		str:      "this is a literal",
		expected: `"this is a literal"`,
	},
	"multiline literal": {
		str: `this is a
literal`,
		expected: `'''this is a
literal'''`,
	},
	"multiline_literal_apostrophe": {
		str: `this is 'a
literal`,
		expected: `"""this is 'a
literal"""`,
	},
	"multiline_literal_quotation": {
		str: `this is "a
literal`,
		expected: `'''this is "a
literal'''`,
	},
}

func TestSanitize(t *testing.T) {
	for name, tc := range sanitizesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := sanitize(tc.str)
			assert.Equal(t, tc.expected, actual, "function should have returned correctly sanitized string")
		})
	}
}
