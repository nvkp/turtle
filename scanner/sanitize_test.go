package scanner

import (
	"testing"

	"github.com/nvkp/turtle/assert"
)

var sanitizeTestCases = map[string]struct {
	base     string
	input    string
	token    string
	label    string
	datatype string
}{
	"with-label": {
		input: `"this is an English text"@en`,
		token: `this is an English text`,
		label: `en`,
	},
	"with-datatype": {
		input:    `"this is an English text"^^xsd:string`,
		token:    `this is an English text`,
		datatype: `xsd:string`,
	},
	"at-in-literal": {
		input:    `"my email is x@y.com"^^xsd:string`,
		token:    `my email is x@y.com`,
		datatype: `xsd:string`,
	},
	"multiline-literal-datatype": {
		input:    `"""Note that SI supports only the use of symbols and deprecates the use of any abbreviations for units."""^^qudt:LatexString`,
		token:    `Note that SI supports only the use of symbols and deprecates the use of any abbreviations for units.`,
		datatype: `qudt:LatexString`,
	},
	"multiline-literal-label": {
		input: `"""Ostrouhej čtyři sta brambor, dokud můžeš ostrou škrabkou!"""@cs`,
		token: `Ostrouhej čtyři sta brambor, dokud můžeš ostrou škrabkou!`,
		label: `cs`,
	},
}

func TestSanitize(t *testing.T) {
	for name, tc := range sanitizeTestCases {
		t.Run(name, func(t *testing.T) {
			s := &Scanner{
				base: tc.base,
			}
			token, label, datatype := s.sanitize(tc.input)
			assert.Equal(t, tc.token, token, "function should have returned correctly sanitized token")
			assert.Equal(t, tc.label, label, "function should have returned correctly extracted label")
			assert.Equal(t, tc.datatype, datatype, "function should have returned correctly extracted datatype")
		})
	}
}
