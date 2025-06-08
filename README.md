# turtle

[![Go Reference](https://pkg.go.dev/badge/github.com/nvkp/turtle.svg)](https://pkg.go.dev/github.com/nvkp/turtle)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

This [Golang](https://go.dev/) package serves as a serializer and parser of the [Turtle](https://www.w3.org/TR/turtle/) format used for representing [RDF](https://www.w3.org/RDF/) data. This package covers most features of the format's **version 1.1**.

## Usage

To add this package as a dependency to you Golang module, run:

```shell
go get github.com/nvkp/turtle
```

The API of the package follows the Golang's traditional pattern for serializing and parsing. The serializing operation happens through the `turtle.Marshal(v interface{}) ([]byte, error)` function. The function accepts the to-be-serialized data as an empty interface and returns the byte slice with the result and possible error value.

It is able to handle single struct, struct, a slice, an array or a pointer to all three. The fields of the structs passed to the function have to be annotated by tags defining which of the fields correspond to which part of the RDF triple.

```golang
var triple = struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}{
	Subject:   "http://e.org/person/Mark_Twain",
	Predicate: "http://e.org/relation/author",
	Object:    "http://e.org/books/Huckleberry_Finn",
}

b, err := turtle.Marshal(&triple)
fmt.Println(string(b)) // <http://e.org/person/Mark_Twain> <http://e.org/relation/author> <http://e.org/books/Huckleberry_Finn> .
```

As default the compact version of the Turtle format is used. The resulting Turtle triples are sorted alphabetically first by subjects, then by predicates and then by objects.

```golang
var triple = []struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}{
	{
		Subject:   "http://e.org/person/Mark_Twain",
		Predicate: "http://e.org/relation/author",
		Object:    "http://e.org/books/Huckleberry_Finn",
	},
	{
		Subject:   "http://e.org/person/Mark_Twain",
		Predicate: "http://e.org/relation/author",
		Object:    "http://e.org/books/Tom_Sawyer",
	},
}

b, err := turtle.Marshal(&triple)
fmt.Println(string(b)) // <http://e.org/person/Mark_Twain> <http://e.org/relation/author> <http://e.org/books/Huckleberry_Finn>, <http://e.org/books/Tom_Sawyer> .
```

Parsing happens through the `turtle.Unmarshal(data []byte, v interface{}) error` function which accepts a byte slice of the turtle data and also a target as a pointer to a struct or to a slice/array of structs that have fields annotated by tags `turtle` defining which field of the struct corresponds to which part of the RDF triple.

```golang
var triple = struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}{}

err := turtle.Unmarshal(
	[]byte("<http://e.org/person/Mark_Twain> <http://e.org/relation/author> <http://e.org/books/Huckleberry_Finn> ."),
	&triple,
)
fmt.Println(triple) // {http://e.org/person/Mark_Twain http://e.org/relation/author http://e.org/books/Huckleberry_Finn}
```

The `turtle.Unmarshal` function accepts the compact version of Turtle just as the N-triples version of the format where each row corresponds to a single triple. It reads `@base` and `@prefix` forms and extends the IRIs that are filled in the target structure with them. It ignores Turtle comments, labels and data types. The keyword `a` gets replaced by `http://www.w3.org/1999/02/22-rdf-syntax-ns#type` IRI. The function is able to handle multiline literals, literal floats, blank nodes, blank node lists and RDF collections.

If the `turtle:"base"` struct tag points at a `string` or `turtle:"prefix"` with `map[string]string` is provided, those fields will be filled in with the base and collection of prefixes respectively. This is per-struct and any future pragma encountered will only effect the following triples. These tags are ignored on marshal, in favor of a configured marshaler. See "Config" for more information.

```golang
var triples = []struct {
	Subject   string            `turtle:"subject"`
	Predicate string            `turtle:"predicate"`
	Object    string            `turtle:"object"`
	Prefixes  map[string]string `turtle:"prefix"`
	Base	  string            `turtle:"base"`
}{}

rdf := `
@base <http://e.org/> .
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix rel: <http://www.perceive.net/schemas/relationship/> .

</green-goblin>
	rel:enemyOf </spiderman> ;
	a foaf:Person ;
	foaf:name "Green Goblin" .
`

err := turtle.Unmarshal(
	[]byte(rdf),
	&triples,
)
```

Both `turtle.Marshal(v interface{}) ([]byte, error)` and `turtle.Unmarshal(data []byte, v interface{}) error` functions can handle two optional field tags `datatype` and `label` annotating the object literals. The struct's attributes with those field tags can either be pointers to string or string values.

```golang
var triples = []struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
	Label     string `turtle:"label"`
	DataType  string `turtle:"datatype"`
}{}

rdf := `
@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix rel: <http://www.perceive.net/schemas/relationship/> .

<http://e.org/green-goblin>
	rel:enemyOf <http://e.org/spiderman> ;
	a foaf:Person ;
	foaf:name "Green Goblin"^^xsd:string , "Zelený Goblin"@cs .
`

err := turtle.Unmarshal(
	[]byte(rdf),
	&triples,
)
```

If you want to resolve URLs automatically at parsing time, create a _configured_ parser with the `turtle.Config` struct. The fields are as follows:

- ResolveURLs: dynamically expand or shorten URLs relative to Base and Prefixes
- Base: configure `@base` without providing syntax
- Prefixes: `map[string]string`, configure prefixes without providing syntax

Base and Prefixes operate exactly like if they were included in the document, and any encountered pragma in a parsed document will affect their representation during unmarshaling.

Example:

```go
c := turtle.Config{
    ResolveURLs: true,
    Base:        "https://example.org/",
    Prefixes:    map[string]string{
        "people": "https://example.org/people/types/"
    },
}

triple := struct {
    Subject     string `turtle:"subject"`
    Predicate   string `turtle:"predicate"`
    Object      string `turtle:"object"`
}{
    Subject: "/people/Mark_Twain",
    Predicate: "a",
    Object: "people:author",
}

data, _ := c.Marshal(triple)
// <https://example.org/people/Mark_Twain> <RDF IRI URL> <https://example.org/people/types/author> .
```

For unmarshaling, `@base` and `@prefix` weigh into the behavior of `ResolveURLs`. They will overwrite any configured options before further resolution. To be absolutely sure what base and prefixes you are using, unmarshal them too.

Example:

```go
c := turtle.Config{
    ResolveURLs: true,
    Base:        "https://example.org/",
    Prefixes:    map[string]string{
        "people": "https://example.org/people/types/"
    },
}

triple := struct {
    Base        string            `turtle:"base"`
    Prefixes    map[string]string `turtle:"prefix"`
    Subject     string            `turtle:"subject"`
    Predicate   string            `turtle:"predicate"`
    Object      string            `turtle:"object"`
}{}

doc := `
@base <https://example2.org> .
</people/Mark_Twain> a people:author .
`

c.Unmarshal([]byte(doc), &triple)

// triple.Base == "https://example2.org"
```

## Existing Alternatives

There is at least one Golang package available on Github that lets you parse and serialize Turtle data: [github.com/deiu/rdf2go](https://github.com/deiu/rdf2go). Its API does not comply with the traditional way of parsing and serializing in Golang programs. It defines its own types appearing in the RDF domain as Triple, Graph, etc.

When a user needs a package that would parse and serialize Turtle data, it is fair to suppose that the user has already defined its own RDF data types as triple or graph. In that case for using the above mentioned package, user has to create a logic for converting its triple data types into the package's data types and adding them to the package's graph structure.

More "Golang way" that this package offers is to annotated the user's already defined structures and the package would read these annotations and behave accordingly.

## Benchmarks

This benchmark compares parsing and serializing operations of the [github.com/deiu/rdf2go](https://github.com/deiu/rdf2go) and [github.com/nvkp/turtle](https://github.com/nvkp/turtle) packages. Both serializing operations are performed on seven triples repeatadly. The parsing operations are performed on a sample of around 27 000 triples. Both parsing and serializing operations from the [github.com/nvkp/turtle](https://github.com/nvkp/turtle) are performed quicker and consume less memory.

```

goos: linux
goarch: amd64
pkg: github.com/nvkp/turtletest
cpu: AMD Ryzen 7 PRO 5850U with Radeon Graphics
BenchmarkMarshalTurtle-16 205250 5801 ns/op
BenchmarkMarshalRDF2Go-16 163112 6384 ns/op
BenchmarkUnmarshalTurtle-16 9 123448964 ns/op
BenchmarkUnmarshalRDF2Go-16 4 261741094 ns/op
PASS
ok github.com/nvkp/turtletest 7.738s

```
