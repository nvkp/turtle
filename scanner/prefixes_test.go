package scanner_test

import (
	"testing"

	"github.com/nvkp/turtle/assert"
	"github.com/nvkp/turtle/scanner"
)

var scannedFile = []byte(`
@base <http://example.org/> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix rel: <http://www.perceive.net/schemas/relationship/> .

<#green-goblin>
	rel:enemyOf <#spiderman> ;
	<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> foaf:Person ;
	foaf:name "Green Goblin".
`)

const expectedBase = `http://example.org/`

var expectedPrefixes = map[string]string{
	"foaf": "http://xmlns.com/foaf/0.1/",
	"rel":  "http://www.perceive.net/schemas/relationship/",
}

func TestScannerBaseAndPrefixes(t *testing.T) {
	s := scanner.New(scannedFile)

	for s.Next() {
		_ = s.Triple()
	}

	assert.Equal(t, expectedBase, s.Base(), "the base collected by the scanner is invalid")
	assert.Equal(t, expectedPrefixes, s.Prefixes(), "the prefixes collected by the scanner are invalid")
}
