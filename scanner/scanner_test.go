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

	"blank_node_property_list": {
		data: []byte(`
		@prefix ericFoaf: <http://www.w3.org/People/Eric/ericP-foaf.rdf#> .
		@prefix : <http://xmlns.com/foaf/0.1/> .
		ericFoaf:ericP :givenName "Eric" ;
					  :knows <http://norman.walsh.name/knows/who/dan-brickley> ,
							  [ :mbox <mailto:timbl@w3.org> ] ,
							  <http://getopenid.com/amyvdh> .
				`),
		expectedTokens: []string{
			`@prefix`,
			`ericFoaf:`,
			`<http://www.w3.org/People/Eric/ericP-foaf.rdf#>`,
			`.`,
			`@prefix`,
			`:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`ericFoaf:ericP`,
			`:givenName`,
			`"Eric"`,
			`;`,
			`:knows`,
			`<http://norman.walsh.name/knows/who/dan-brickley>`,
			`,`,
			`[`,
			`:mbox`,
			`<mailto:timbl@w3.org>`,
			`]`,
			`,`,
			`<http://getopenid.com/amyvdh>`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"http://www.w3.org/People/Eric/ericP-foaf.rdf#ericP", "http://xmlns.com/foaf/0.1/givenName", "Eric"},
			{"http://www.w3.org/People/Eric/ericP-foaf.rdf#ericP", "http://xmlns.com/foaf/0.1/knows", "http://norman.walsh.name/knows/who/dan-brickley"},
			{"_:b0", "http://xmlns.com/foaf/0.1/mbox", "mailto:timbl@w3.org"},
			{"http://www.w3.org/People/Eric/ericP-foaf.rdf#ericP", "http://xmlns.com/foaf/0.1/knows", "_:b0"},
			{"http://www.w3.org/People/Eric/ericP-foaf.rdf#ericP", "http://xmlns.com/foaf/0.1/knows", "http://getopenid.com/amyvdh"},
		},
	},
	"blank_node_property_list_nested": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .

		foaf:Alice foaf:knows [
			foaf:name "Bob" ;
			foaf:knows [
				foaf:name "Eve" ] ;
			foaf:mbox <bob@example.com> ] .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`foaf:Alice`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Bob"`,
			`;`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Eve"`,
			`]`,
			`;`,
			`foaf:mbox`,
			`<bob@example.com>`,
			`]`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/name", "Bob"},
			{"_:b1", "http://xmlns.com/foaf/0.1/name", "Eve"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "_:b1"},
			{"_:b0", "http://xmlns.com/foaf/0.1/mbox", "bob@example.com"},
			{"http://xmlns.com/foaf/0.1/Alice", "http://xmlns.com/foaf/0.1/knows", "_:b0"},
		},
	},
	"blank_node_as_subject": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
		# Someone knows someone else, who has the name "Bob".

		[ foaf:name "Bob" ] foaf:knows foaf:someone .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`[`,
			`foaf:name`,
			`"Bob"`,
			`]`,
			`foaf:knows`,
			`foaf:someone`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/name", "Bob"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "http://xmlns.com/foaf/0.1/someone"},
		},
	},
	"blank_node_as_subject_nested": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
		# Someone knows someone else, who has the name "Bob".

		[ foaf:name "Bob"; foaf:knows [ foaf:name "Alice" ] ] foaf:knows foaf:someone .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`[`,
			`foaf:name`,
			`"Bob"`,
			`;`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Alice"`,
			`]`,
			`]`,
			`foaf:knows`,
			`foaf:someone`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/name", "Bob"},
			{"_:b1", "http://xmlns.com/foaf/0.1/name", "Alice"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "_:b1"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "http://xmlns.com/foaf/0.1/someone"},
		},
	},
	"blank_node_empty_subject": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .

		# Someone knows someone else, who has the name "Bob".
		[ ] foaf:knows foaf:someone .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`[`,
			`]`,
			`foaf:knows`,
			`foaf:someone`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "http://xmlns.com/foaf/0.1/someone"},
		},
	},
	"blank_node_empty_subject_no_whitespace": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .

		# Someone knows someone else, who has the name "Bob".
		[] foaf:knows foaf:someone .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`[`,
			`]`,
			`foaf:knows`,
			`foaf:someone`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "http://xmlns.com/foaf/0.1/someone"},
		},
	},
	"blank_node_collision_avoided": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .

		_:b0 foaf:knows [
			foaf:name "Bob" ;
			foaf:knows [
				foaf:name "Eve" ] ;
			foaf:mbox <bob@example.com> ] .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`_:b0`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Bob"`,
			`;`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Eve"`,
			`]`,
			`;`,
			`foaf:mbox`,
			`<bob@example.com>`,
			`]`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b1", "http://xmlns.com/foaf/0.1/name", "Bob"},
			{"_:b2", "http://xmlns.com/foaf/0.1/name", "Eve"},
			{"_:b1", "http://xmlns.com/foaf/0.1/knows", "_:b2"},
			{"_:b1", "http://xmlns.com/foaf/0.1/mbox", "bob@example.com"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "_:b1"},
		},
	},
	"blank_node_subject_and_object": {
		data: []byte(`
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
		# Someone knows someone else, who has the name "Bob".
		[ foaf:name "Bob" ] foaf:knows [ foaf:name "Alice" ] .
				`),
		expectedTokens: []string{
			`@prefix`,
			`foaf:`,
			`<http://xmlns.com/foaf/0.1/>`,
			`.`,
			`[`,
			`foaf:name`,
			`"Bob"`,
			`]`,
			`foaf:knows`,
			`[`,
			`foaf:name`,
			`"Alice"`,
			`]`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://xmlns.com/foaf/0.1/name", "Bob"},
			{"_:b1", "http://xmlns.com/foaf/0.1/name", "Alice"},
			{"_:b0", "http://xmlns.com/foaf/0.1/knows", "_:b1"},
		},
	},
	"collection_object": {
		data: []byte(`
		@prefix : <http://example.org/stuff/1.0/> .
		:a :b ( "apple" "banana" ) .
				`),
		expectedTokens: []string{
			`@prefix`,
			`:`,
			`<http://example.org/stuff/1.0/>`,
			`.`,
			`:a`,
			`:b`,
			`(`,
			`"apple"`,
			`"banana"`,
			`)`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `apple`},
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "_:b1"},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `banana`},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"},
			{"http://example.org/stuff/1.0/a", "http://example.org/stuff/1.0/b", "_:b0"},
		},
	},
	"collection_object_empty": {
		data: []byte(`
		@prefix : <http://example.org/stuff/1.0/> .
		:subject :predicate2 () .
				`),
		expectedTokens: []string{
			`@prefix`,
			`:`,
			`<http://example.org/stuff/1.0/>`,
			`.`,
			`:subject`,
			`:predicate2`,
			`(`,
			`)`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"http://example.org/stuff/1.0/subject", "http://example.org/stuff/1.0/predicate2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"},
		},
	},
	"collection_subject": {
		data: []byte(`
		@prefix : <http://example.org/stuff/1.0/> .
		(1 2.0 3E1) :p "w" .
				`),
		expectedTokens: []string{
			`@prefix`,
			`:`,
			`<http://example.org/stuff/1.0/>`,
			`.`,
			`(`,
			`1`,
			`2.0`,
			`3E1`,
			`)`,
			`:p`,
			`"w"`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `1`},
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "_:b1"},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `2.0`},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "_:b2"},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `3E1`},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"},
			{"_:b0", "http://example.org/stuff/1.0/p", "w"},
		},
	},
	"collection_sanitized": {
		data: []byte(`
		@prefix : <http://example.org/stuff/1.0/> .
		:a :b ( "apple"@en :c ) .
				`),
		expectedTokens: []string{
			`@prefix`,
			`:`,
			`<http://example.org/stuff/1.0/>`,
			`.`,
			`:a`,
			`:b`,
			`(`,
			`"apple"@en`,
			`:c`,
			`)`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `apple`},
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "_:b1"},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `http://example.org/stuff/1.0/c`},
			{"_:b1", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"},
			{"http://example.org/stuff/1.0/a", "http://example.org/stuff/1.0/b", "_:b0"},
		},
	},
	"collection_nested": {
		data: []byte(`
		@prefix : <http://example.org/stuff/1.0/> .
		(1 [:p :q] ( 2 ) ) :p2 :q2 .
				`),
		expectedTokens: []string{
			`@prefix`,
			`:`,
			`<http://example.org/stuff/1.0/>`,
			`.`,
			`(`,
			`1`,
			`[`,
			`:p`,
			`:q`,
			`]`,
			`(`,
			`2`,
			`)`,
			`)`,
			`:p2`,
			`:q2`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b1", "http://example.org/stuff/1.0/p", `http://example.org/stuff/1.0/q`},
			{"_:b3", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `2`},
			{"_:b3", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", `http://www.w3.org/1999/02/22-rdf-syntax-ns#nil`},
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `1`},
			{"_:b0", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", `_:b2`},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `_:b1`},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", `_:b4`},
			{"_:b4", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", `_:b3`},
			{"_:b4", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", `http://www.w3.org/1999/02/22-rdf-syntax-ns#nil`},
			{"_:b0", "http://example.org/stuff/1.0/p2", `http://example.org/stuff/1.0/q2`},
		},
	},
	"blank_node_list_in_collection_in_blank_node_list": {
		data: []byte(`
		@prefix brick: <https://brickschema.org/schema/Brick#> .
		@prefix sh: <http://www.w3.org/ns/shacl#> .
		@prefix tag: <https://brickschema.org/schema/BrickTag#> .

		brick:Portfolio sh:property [ sh:or ( [ sh:class brick:Site ] ) ;
				sh:path brick:hasPart ] ;
		sh:rule [ a sh:TripleRule ;
				sh:object tag:Collection ;
				sh:predicate brick:hasTag ;
				sh:subject sh:this ],
			[ a sh:TripleRule ;
				sh:object tag:Portfolio ;
				sh:predicate brick:hasTag ;
				sh:subject sh:this ] ;
		brick:hasAssociatedTag tag:Collection,
			tag:Portfolio .
				`),
		expectedTokens: []string{
			`@prefix`,
			`brick:`,
			`<https://brickschema.org/schema/Brick#>`,
			`.`,
			`@prefix`,
			`sh:`,
			`<http://www.w3.org/ns/shacl#>`,
			`.`,
			`@prefix`,
			`tag:`,
			`<https://brickschema.org/schema/BrickTag#>`,
			`.`,
			`brick:Portfolio`,
			`sh:property`,
			`[`,
			`sh:or`,
			`(`,
			`[`,
			`sh:class`,
			`brick:Site`,
			`]`,
			`)`,
			`;`,
			`sh:path`,
			`brick:hasPart`,
			`]`,
			`;`,
			`sh:rule`,
			`[`,
			`a`,
			`sh:TripleRule`,
			`;`,
			`sh:object`,
			`tag:Collection`,
			`;`,
			`sh:predicate`,
			`brick:hasTag`,
			`;`,
			`sh:subject`,
			`sh:this`,
			`]`,
			`,`,
			`[`,
			`a`,
			`sh:TripleRule`,
			`;`,
			`sh:object`,
			`tag:Portfolio`,
			`;`,
			`sh:predicate`,
			`brick:hasTag`,
			`;`,
			`sh:subject`,
			`sh:this`,
			`]`,
			`;`,
			`brick:hasAssociatedTag`,
			`tag:Collection`,
			`,`,
			`tag:Portfolio`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"_:b1", "http://www.w3.org/ns/shacl#class", `https://brickschema.org/schema/Brick#Site`},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#first", "_:b1"},
			{"_:b2", "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest", "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"},
			{"_:b0", "http://www.w3.org/ns/shacl#or", "_:b2"},
			{"_:b0", "http://www.w3.org/ns/shacl#path", "https://brickschema.org/schema/Brick#hasPart"},
			{"https://brickschema.org/schema/Brick#Portfolio", "http://www.w3.org/ns/shacl#property", "_:b0"},
			{"_:b3", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/ns/shacl#TripleRule"},
			{"_:b3", "http://www.w3.org/ns/shacl#object", "https://brickschema.org/schema/BrickTag#Collection"},
			{"_:b3", "http://www.w3.org/ns/shacl#predicate", "https://brickschema.org/schema/Brick#hasTag"},
			{"_:b3", "http://www.w3.org/ns/shacl#subject", "http://www.w3.org/ns/shacl#this"},
			{"https://brickschema.org/schema/Brick#Portfolio", "http://www.w3.org/ns/shacl#rule", "_:b3"},
			{"_:b4", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", "http://www.w3.org/ns/shacl#TripleRule"},
			{"_:b4", "http://www.w3.org/ns/shacl#object", "https://brickschema.org/schema/BrickTag#Portfolio"},
			{"_:b4", "http://www.w3.org/ns/shacl#predicate", "https://brickschema.org/schema/Brick#hasTag"},
			{"_:b4", "http://www.w3.org/ns/shacl#subject", "http://www.w3.org/ns/shacl#this"},
			{"https://brickschema.org/schema/Brick#Portfolio", "http://www.w3.org/ns/shacl#rule", "_:b4"},
			{"https://brickschema.org/schema/Brick#Portfolio", "https://brickschema.org/schema/Brick#hasAssociatedTag", "https://brickschema.org/schema/BrickTag#Collection"},
			{"https://brickschema.org/schema/Brick#Portfolio", "https://brickschema.org/schema/Brick#hasAssociatedTag", "https://brickschema.org/schema/BrickTag#Portfolio"},
		},
	},
	"literal_character_in_literal": {
		data: []byte(`
		@prefix unit: <http://qudt.org/vocab/unit/> .
		@prefix qudt: <http://qudt.org/schema/qudt/> .
		@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
		unit:ARCMIN a qudt:Unit ;
		rdfs:label "ArcMinute"@en ;
		qudt:symbol "'",
			"'"^^xsd:string .
				`),
		expectedTokens: []string{
			`@prefix`,
			`unit:`,
			`<http://qudt.org/vocab/unit/>`,
			`.`,
			`@prefix`,
			`qudt:`,
			`<http://qudt.org/schema/qudt/>`,
			`.`,
			`@prefix`,
			`rdfs:`,
			`<http://www.w3.org/2000/01/rdf-schema#>`,
			`.`,
			`unit:ARCMIN`,
			`a`,
			`qudt:Unit`,
			`;`,
			`rdfs:label`,
			`"ArcMinute"@en`,
			`;`,
			`qudt:symbol`,
			`"'"`,
			`,`,
			`"'"^^xsd:string`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"http://qudt.org/vocab/unit/ARCMIN", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", `http://qudt.org/schema/qudt/Unit`},
			{"http://qudt.org/vocab/unit/ARCMIN", "http://www.w3.org/2000/01/rdf-schema#label", `ArcMinute`},
			{"http://qudt.org/vocab/unit/ARCMIN", "http://qudt.org/schema/qudt/symbol", `'`},
			{"http://qudt.org/vocab/unit/ARCMIN", "http://qudt.org/schema/qudt/symbol", `'`},
		},
	},
	"float_in_iri": {
		data: []byte(`
		@prefix brick: <https://brickschema.org/schema/Brick#> .
		brick:PM2.5_Sensor brick:hasQuantity brick:PM2.5_Concentration .
				`),
		expectedTokens: []string{
			`@prefix`,
			`brick:`,
			`<https://brickschema.org/schema/Brick#>`,
			`.`,
			`brick:PM2.5_Sensor`,
			`brick:hasQuantity`,
			`brick:PM2.5_Concentration`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"https://brickschema.org/schema/Brick#PM2.5_Sensor", "https://brickschema.org/schema/Brick#hasQuantity", `https://brickschema.org/schema/Brick#PM2.5_Concentration`},
		},
	},
	"float-with-exponents": {
		data: []byte(`
		@prefix unit: <http://qudt.org/vocab/unit/> .
		@prefix qudt: <http://qudt.org/schema/qudt/> .
		@prefix quantitykind: <http://qudt.org/vocab/quantitykind/> .
		unit:A
		a qudt:Unit ;
		qudt:conversionMultiplierSN 1.0E0, 42E3, 1e0, -2.3E-12, +.3e+2 ;
		qudt:hasQuantityKind quantitykind:TotalCurrent .
				`),
		expectedTokens: []string{
			`@prefix`,
			`unit:`,
			`<http://qudt.org/vocab/unit/>`,
			`.`,
			`@prefix`,
			`qudt:`,
			`<http://qudt.org/schema/qudt/>`,
			`.`,
			`@prefix`,
			`quantitykind:`,
			`<http://qudt.org/vocab/quantitykind/>`,
			`.`,
			`unit:A`,
			`a`,
			`qudt:Unit`,
			`;`,
			`qudt:conversionMultiplierSN`,
			`1.0E0`,
			`,`,
			`42E3`,
			`,`,
			`1e0`,
			`,`,
			`-2.3E-12`,
			`,`,
			`+.3e+2`,
			`;`,
			`qudt:hasQuantityKind`,
			`quantitykind:TotalCurrent`,
			`.`,
		},
		expectedTriples: [][3]string{
			{"http://qudt.org/vocab/unit/A", "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", `http://qudt.org/schema/qudt/Unit`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/conversionMultiplierSN", `1.0E0`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/conversionMultiplierSN", `42E3`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/conversionMultiplierSN", `1e0`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/conversionMultiplierSN", `-2.3E-12`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/conversionMultiplierSN", `+.3e+2`},
			{"http://qudt.org/vocab/unit/A", "http://qudt.org/schema/qudt/hasQuantityKind", `http://qudt.org/vocab/quantitykind/TotalCurrent`},
		},
	},
}

func TestScanTurtle(t *testing.T) {
	for name, tc := range scanTestCases {
		t.Run(name, func(t *testing.T) {
			s := bufio.NewScanner(bytes.NewReader(tc.data))
			s.Split(splitTurtle)
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
