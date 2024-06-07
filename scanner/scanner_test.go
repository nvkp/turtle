package scanner

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/nvkp/turtle/assert"
)

var scanTestCases = map[string]struct {
	data            []byte
	expectedTokens  []string
	expectedTriples [][3]string
}{
	"spiderman compact": {
		data: []byte(`<http://example.org/green-goblin>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
					<http://xmlns.com/foaf/0.1/name> "Green Goblin".<http://example.org/spiderman>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin>;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person>;
					<http://xmlns.com/foaf/0.1/name> "Spiderman", "Человек-паук" .`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"`,
			".",
			"<http://example.org/spiderman>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/green-goblin>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Spiderman"`,
			",",
			`"Человек-паук"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
			{"http://example.org/spiderman", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/green-goblin"},
			{"http://example.org/spiderman", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Spiderman"},
			{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Человек-паук"},
		},
	},
	"ignore_comments": {
		data: []byte(`<http://example.org/green-goblin>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ; # this is a comment
					<http://xmlns.com/foaf/0.1/name> "Green Goblin".`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"ignore_label": {
		data: []byte(`<http://example.org/green-goblin>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
					<http://xmlns.com/foaf/0.1/name> "Green Goblin"@en .`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"@en`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"ignore_prefixed_datatype": {
		data: []byte(`<http://example.org/green-goblin>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
					<http://xmlns.com/foaf/0.1/name> "Green Goblin"^^xsd:string .`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"^^xsd:string`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"booleans": {
		data: []byte(`@prefix s: <http://example.org/stats/> .
					<http://somecountry.example/census2007>
						s:isLandlocked false .`),
		expectedTokens: []string{
			"@prefix",
			"s:",
			"<http://example.org/stats/>",
			".",
			"<http://somecountry.example/census2007>",
			"s:isLandlocked",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats/isLandlocked", "false"},
		},
	},
	"empty_prefix": {
		data: []byte(`@prefix : <http://example.org/stats/> .
					<http://somecountry.example/census2007>
						:isLandlocked false .`),
		expectedTokens: []string{
			"@prefix",
			":",
			"<http://example.org/stats/>",
			".",
			"<http://somecountry.example/census2007>",
			":isLandlocked",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats/isLandlocked", "false"},
		},
	},
	"prefix_no_ending_slash": {
		data: []byte(`@prefix : <http://example.org/stats> .
					<http://somecountry.example/census2007>
						:isLandlocked false .`),
		expectedTokens: []string{
			"@prefix",
			":",
			"<http://example.org/stats>",
			".",
			"<http://somecountry.example/census2007>",
			":isLandlocked",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats/isLandlocked", "false"},
		},
	},
	"base_no_ending_slash": {
		data: []byte(`@base <http://example.org/stats> .
					<http://somecountry.example/census2007>
						<#isLandlocked> false .`),
		expectedTokens: []string{
			"@base",
			"<http://example.org/stats>",
			".",
			"<http://somecountry.example/census2007>",
			"<#isLandlocked>",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats/isLandlocked", "false"},
		},
	},
	"base_with_ending_slash": {
		data: []byte(`@base <http://example.org/stats/> .
					<http://somecountry.example/census2007>
						<#isLandlocked> false .`),
		expectedTokens: []string{
			"@base",
			"<http://example.org/stats/>",
			".",
			"<http://somecountry.example/census2007>",
			"<#isLandlocked>",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats/isLandlocked", "false"},
		},
	},
	"ignore_datatype": {
		data: []byte(`<http://example.org/green-goblin>
					<http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> ;
					<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> ;
					<http://xmlns.com/foaf/0.1/name> "Green Goblin"^^<http://www.w3.org/2001/XMLSchema#string> .`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			";",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"^^<http://www.w3.org/2001/XMLSchema#string>`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"read_prefix": {
		data: []byte(`@prefix foaf: <http://xmlns.com/foaf/0.1/> .
					@prefix rel: <http://www.perceive.net/schemas/relationship/> .

					<http://example.org/green-goblin>
						rel:enemyOf <http://example.org/spiderman> ;
						<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> foaf:Person ;
						foaf:name "Green Goblin".`),
		expectedTokens: []string{
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"@prefix",
			"rel:",
			"<http://www.perceive.net/schemas/relationship/>",
			".",
			"<http://example.org/green-goblin>",
			"rel:enemyOf",
			"<http://example.org/spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"foaf:Person",
			";",
			"foaf:name",
			`"Green Goblin"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"read_prefix_and_base": {
		data: []byte(`@base <http://example.org/> .
					@prefix foaf: <http://xmlns.com/foaf/0.1/> .
					@prefix rel: <http://www.perceive.net/schemas/relationship/> .

					<#green-goblin>
						rel:enemyOf <#spiderman> ;
						<http://www.w3.org/1999/02/22-rdf-syntax-ns#type> foaf:Person ;
						foaf:name "Green Goblin".`),
		expectedTokens: []string{
			"@base",
			"<http://example.org/>",
			".",
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"@prefix",
			"rel:",
			"<http://www.perceive.net/schemas/relationship/>",
			".",
			"<#green-goblin>",
			"rel:enemyOf",
			"<#spiderman>",
			";",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"foaf:Person",
			";",
			"foaf:name",
			`"Green Goblin"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"spiderman n-triples": {
		data: []byte(`<http://example.org/green-goblin> <http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/spiderman> .
					<http://example.org/green-goblin> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
					<http://example.org/green-goblin> <http://xmlns.com/foaf/0.1/name> "Green Goblin".
					<http://example.org/spiderman> <http://www.perceive.net/schemas/relationship/enemyOf> <http://example.org/green-goblin> .
					<http://example.org/spiderman> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
					<http://example.org/spiderman> <http://xmlns.com/foaf/0.1/name> "Spiderman" .
					<http://example.org/spiderman> <http://xmlns.com/foaf/0.1/name> "Человек-паук" .`),
		expectedTokens: []string{
			"<http://example.org/green-goblin>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/spiderman>",
			".",
			"<http://example.org/green-goblin>",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			".",
			"<http://example.org/green-goblin>",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Green Goblin"`,
			".",
			"<http://example.org/spiderman>",
			"<http://www.perceive.net/schemas/relationship/enemyOf>",
			"<http://example.org/green-goblin>",
			".",
			"<http://example.org/spiderman>",
			"<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>",
			"<http://xmlns.com/foaf/0.1/Person>",
			".",
			"<http://example.org/spiderman>",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Spiderman"`,
			".",
			"<http://example.org/spiderman>",
			"<http://xmlns.com/foaf/0.1/name>",
			`"Человек-паук"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
			{"http://example.org/spiderman", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/green-goblin"},
			{"http://example.org/spiderman", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Spiderman"},
			{"http://example.org/spiderman", "http://xmlns.com/foaf/0.1/name", "Человек-паук"},
		},
	},
	"read_rdf_type_shorthand": {
		data: []byte(`@prefix foaf: <http://xmlns.com/foaf/0.1/> .
					@prefix rel: <http://www.perceive.net/schemas/relationship/> .

					<http://example.org/green-goblin>
						rel:enemyOf <http://example.org/spiderman> ;
						a foaf:Person ;
						foaf:name "Green Goblin".`),
		expectedTokens: []string{
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"@prefix",
			"rel:",
			"<http://www.perceive.net/schemas/relationship/>",
			".",
			"<http://example.org/green-goblin>",
			"rel:enemyOf",
			"<http://example.org/spiderman>",
			";",
			"a",
			"foaf:Person",
			";",
			"foaf:name",
			`"Green Goblin"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://www.perceive.net/schemas/relationship/enemyOf", "http://example.org/spiderman"},
			{"http://example.org/green-goblin", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://xmlns.com/foaf/0.1/Person"},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "Green Goblin"},
		},
	},
	"apostrophe_literal": {
		data: []byte(`@prefix foaf: <http://xmlns.com/foaf/0.1/> .

					<http://example.org/green-goblin> foaf:name 'Weird Name With " in it' .`),
		expectedTokens: []string{
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"<http://example.org/green-goblin>",
			"foaf:name",
			`'Weird Name With " in it'`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", `Weird Name With " in it`},
		},
	},
	"apostrophe_in_quotation_mark_literal": {
		data: []byte(`@prefix foaf: <http://xmlns.com/foaf/0.1/> .

					<http://example.org/green-goblin> foaf:name "Weird Name With ' in it" .`),
		expectedTokens: []string{
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"<http://example.org/green-goblin>",
			"foaf:name",
			`"Weird Name With ' in it"`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", `Weird Name With ' in it`},
		},
	},
	"mind_gt_lt_in_literal": {
		data: []byte(`@prefix foaf: <http://xmlns.com/foaf/0.1/> .

					<http://example.org/green-goblin> foaf:name "Weird Name With < and > and < in it", <http://example.org/some-iri> .`),
		expectedTokens: []string{
			"@prefix",
			"foaf:",
			"<http://xmlns.com/foaf/0.1/>",
			".",
			"<http://example.org/green-goblin>",
			"foaf:name",
			`"Weird Name With < and > and < in it"`,
			",",
			`<http://example.org/some-iri>`,
			".",
		},
		expectedTriples: [][3]string{
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", `Weird Name With < and > and < in it`},
			{"http://example.org/green-goblin", "http://xmlns.com/foaf/0.1/name", "http://example.org/some-iri"},
		},
	},
	"quation-mark-multiline-literal": {
		data: []byte(`@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
				@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
				@prefix schema: <https://schema.org/> .

				schema:ComicSeries a rdfs:Class ;
					rdfs:label "ComicSeries" ;
					rdfs:comment """A sequential publication of comic stories under a
		unifying title, for example "The Amazing Spider-Man" or "Groo the
		Wanderer".""" ;
					rdfs:subClassOf schema:Periodical ;
					schema:isPartOf <https://bib.schema.org> .`),
		expectedTokens: []string{
			`@prefix`,
			`rdf:`,
			`<http://www.w3.org/1999/02/22-rdf-syntax-ns#>`,
			`.`,
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`@prefix`,
			`schema:`,
			`<https://schema.org/>`,
			`.`,
			`schema:ComicSeries`,
			`a`,
			`rdfs:Class`,
			`;`,
			`rdfs:label`,
			`"ComicSeries"`,
			`;`,
			`rdfs:comment`,
			`"""A sequential publication of comic stories under a
		unifying title, for example "The Amazing Spider-Man" or "Groo the
		Wanderer"."""`,
			`;`,
			`rdfs:subClassOf`,
			`schema:Periodical`,
			`;`,
			`schema:isPartOf`,
			`<https://bib.schema.org>`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://schema.org/ComicSeries", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/2000/01/rdf-schema#Class"},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#label", "ComicSeries"},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#comment", `A sequential publication of comic stories under a
		unifying title, for example "The Amazing Spider-Man" or "Groo the
		Wanderer".`},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#subClassOf", "https://schema.org/Periodical"},
			{"https://schema.org/ComicSeries", "https://schema.org/isPartOf", "https://bib.schema.org"},
		},
	},
	"apostrophe-multiline-literal": {
		data: []byte(`@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
					@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
					@prefix schema: <https://schema.org/> .

					schema:ComicSeries a rdfs:Class ;
						rdfs:label "ComicSeries" ;
						rdfs:comment '''A sequential publication of comic stories under a
			unifying title, for example "The Amazing Spider-Man" or "Groo the
			Wanderer".''' ;
						rdfs:subClassOf schema:Periodical ;
						schema:isPartOf <https://bib.schema.org> .`),
		expectedTokens: []string{
			`@prefix`,
			`rdf:`,
			`<http://www.w3.org/1999/02/22-rdf-syntax-ns#>`,
			`.`,
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`@prefix`,
			`schema:`,
			`<https://schema.org/>`,
			`.`,
			`schema:ComicSeries`,
			`a`,
			`rdfs:Class`,
			`;`,
			`rdfs:label`,
			`"ComicSeries"`,
			`;`,
			`rdfs:comment`,
			`'''A sequential publication of comic stories under a
			unifying title, for example "The Amazing Spider-Man" or "Groo the
			Wanderer".'''`,
			`;`,
			`rdfs:subClassOf`,
			`schema:Periodical`,
			`;`,
			`schema:isPartOf`,
			`<https://bib.schema.org>`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://schema.org/ComicSeries", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/2000/01/rdf-schema#Class"},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#label", "ComicSeries"},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#comment", `A sequential publication of comic stories under a
			unifying title, for example "The Amazing Spider-Man" or "Groo the
			Wanderer".`},
			{"https://schema.org/ComicSeries", "http://www.w3.org/2000/01/rdf-schema#subClassOf", "https://schema.org/Periodical"},
			{"https://schema.org/ComicSeries", "https://schema.org/isPartOf", "https://bib.schema.org"},
		},
	},
	"escaped-quation": {
		data: []byte(`
				@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
				@prefix schema: <https://schema.org/> .
				schema:FAQPage a rdfs:Class ;
		rdfs:label "FAQPage" ;
		rdfs:comment "A [[FAQPage]] is a [[WebPage]] presenting one or more \"[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\" (see also [[QAPage]])." ;
		rdfs:subClassOf schema:WebPage ;
		schema:source <https://github.com/schemaorg/schemaorg/issues/1723> .`),
		expectedTokens: []string{
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`@prefix`,
			`schema:`,
			`<https://schema.org/>`,
			`.`,
			`schema:FAQPage`,
			`a`,
			`rdfs:Class`,
			`;`,
			`rdfs:label`,
			`"FAQPage"`,
			`;`,
			`rdfs:comment`,
			`"A [[FAQPage]] is a [[WebPage]] presenting one or more \"[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\" (see also [[QAPage]])."`,
			`;`,
			`rdfs:subClassOf`,
			`schema:WebPage`,
			`;`,
			`schema:source`,
			`<https://github.com/schemaorg/schemaorg/issues/1723>`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://schema.org/FAQPage", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/2000/01/rdf-schema#Class"},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#label", "FAQPage"},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#comment", `A [[FAQPage]] is a [[WebPage]] presenting one or more \"[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\" (see also [[QAPage]]).`},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#subClassOf", "https://schema.org/WebPage"},
			{"https://schema.org/FAQPage", "https://schema.org/source", "https://github.com/schemaorg/schemaorg/issues/1723"},
		},
	},
	"escaped-apostrophe": {
		data: []byte(`
			@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
			@prefix schema: <https://schema.org/> .
			schema:FAQPage a rdfs:Class ;
	rdfs:label "FAQPage" ;
	rdfs:comment 'A [[FAQPage]] is a [[WebPage]] presenting one or more \'[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\' (see also [[QAPage]]).' ;
	rdfs:subClassOf schema:WebPage ;
	schema:source <https://github.com/schemaorg/schemaorg/issues/1723> .`),
		expectedTokens: []string{
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`@prefix`,
			`schema:`,
			`<https://schema.org/>`,
			`.`,
			`schema:FAQPage`,
			`a`,
			`rdfs:Class`,
			`;`,
			`rdfs:label`,
			`"FAQPage"`,
			`;`,
			`rdfs:comment`,
			`'A [[FAQPage]] is a [[WebPage]] presenting one or more \'[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\' (see also [[QAPage]]).'`,
			`;`,
			`rdfs:subClassOf`,
			`schema:WebPage`,
			`;`,
			`schema:source`,
			`<https://github.com/schemaorg/schemaorg/issues/1723>`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://schema.org/FAQPage", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/2000/01/rdf-schema#Class"},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#label", "FAQPage"},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#comment", `A [[FAQPage]] is a [[WebPage]] presenting one or more \'[Frequently asked questions](https://en.wikipedia.org/wiki/FAQ)\' (see also [[QAPage]]).`},
			{"https://schema.org/FAQPage", "http://www.w3.org/2000/01/rdf-schema#subClassOf", "https://schema.org/WebPage"},
			{"https://schema.org/FAQPage", "https://schema.org/source", "https://github.com/schemaorg/schemaorg/issues/1723"},
		},
	},
	"base_with_number_sign": {
		data: []byte(`@base <http://example.org/stats#> .
					<http://somecountry.example/census2007>
						<#isLandlocked> false .`),
		expectedTokens: []string{
			"@base",
			"<http://example.org/stats#>",
			".",
			"<http://somecountry.example/census2007>",
			"<#isLandlocked>",
			"false",
			".",
		},
		expectedTriples: [][3]string{
			{"http://somecountry.example/census2007", "http://example.org/stats#isLandlocked", "false"},
		},
	},
	"prefix_with_number_sign": {
		data: []byte(`
@prefix dcterms: <http://purl.org/dc/terms/> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix schema: <https://schema.org/> .

schema:identifier a rdf:Property ;
    rdfs:label "identifier" ;
    owl:equivalentProperty dcterms:identifier ;
    schema:domainIncludes schema:Thing ;
    schema:rangeIncludes schema:PropertyValue,
        schema:Text,
        schema:URL .
		`),
		expectedTokens: []string{
			`@prefix`,
			`dcterms:`,
			`<http://purl.org/dc/terms/>`,
			`.`,
			`@prefix`,
			`owl:`,
			`<http://www.w3.org/2002/07/owl#>`,
			`.`,
			`@prefix`,
			`rdf:`,
			`<http://www.w3.org/1999/02/22-rdf-syntax-ns#>`,
			`.`,
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`@prefix`,
			`schema:`,
			`<https://schema.org/>`,
			`.`,
			`schema:identifier`,
			`a`,
			`rdf:Property`,
			`;`,
			`rdfs:label`,
			`"identifier"`,
			`;`,
			`owl:equivalentProperty`,
			`dcterms:identifier`,
			`;`,
			`schema:domainIncludes`,
			`schema:Thing`,
			`;`,
			`schema:rangeIncludes`,
			`schema:PropertyValue`,
			`,`,
			`schema:Text`,
			`,`,
			`schema:URL`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://schema.org/identifier", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"},
			{"https://schema.org/identifier", "http://www.w3.org/2000/01/rdf-schema#label", "identifier"},
			{"https://schema.org/identifier", "http://www.w3.org/2002/07/owl#equivalentProperty", "http://purl.org/dc/terms/identifier"},
			{"https://schema.org/identifier", "https://schema.org/domainIncludes", "https://schema.org/Thing"},
			{"https://schema.org/identifier", "https://schema.org/rangeIncludes", "https://schema.org/PropertyValue"},
			{"https://schema.org/identifier", "https://schema.org/rangeIncludes", "https://schema.org/Text"},
			{"https://schema.org/identifier", "https://schema.org/rangeIncludes", "https://schema.org/URL"},
		},
	},
}

func TestScanTurtle(t *testing.T) {
	for name, tc := range scanTestCases {
		t.Run(name, func(t *testing.T) {
			s := bufio.NewScanner(bytes.NewReader(tc.data))
			s.Split(scanTurtle)
			actual := make([]string, 0)
			for {
				ok := s.Scan()
				if !ok {
					break
				}
				actual = append(actual, s.Text())
			}
			assert.Equal(t, tc.expectedTokens, actual, "scanTurtle should have created correct turtle tokens")
		})
	}
}

func TestNext(t *testing.T) {
	for name, tc := range scanTestCases {
		t.Run(name, func(t *testing.T) {
			s := New(tc.data)
			actual := make([][3]string, 0)
			for {
				ok := s.Next()
				if !ok {
					break
				}
				actual = append(actual, s.Triple())
			}
			assert.Equal(t, tc.expectedTriples, actual, "scanner should have created correct turtle triples")
		})
	}
}
