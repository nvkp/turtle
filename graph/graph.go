package graph

import (
	"fmt"
	"sort"
)

// Options changes the behavior of the graph. It is passed to NewWithOptions.
type Options struct {
	// If set, will output a `@base` pragma at the start. Will normalize all URLs
	// that start with the base to their relative components.
	Base string
	// If set, any encountering of the prefix URL prefixes will be normalized to
	// use the prefix. Additionally, @prefix lines are output at the top of the
	// document for each one.
	Prefixes map[string]string
}

type object struct {
	item     string
	typ      string
	datatype string
	label    string
}

// Graph serves as a buffer that consumes triples one by one
// and can return a byte slice containing Turtle data of all
// triples consumed.
type Graph struct {
	options Options
	m       map[string]map[string][]object
}

// New returns a pointer to a new instance of graph.Graph. No options are set.
func New() *Graph {
	return NewWithOptions(Options{})
}

// NewWithOptions constructs a graph with options to tweak its behavior. See Options.
func NewWithOptions(options Options) *Graph {
	return &Graph{
		options: options,
		m:       make(map[string]map[string][]object),
	}
}

// Accept stores a new triple to the graph.
func (g *Graph) Accept(t [3]string) error {
	return g.accept(t[0], t[1], object{item: t[2]})
}

// AcceptWithAnnotations stores a new triple with eventual label and data type
// of the object literal to the graph.
func (g *Graph) AcceptWithAnnotations(t [6]string) error {
	return g.accept(t[0], t[1], object{item: t[2], label: t[3], datatype: t[4], typ: t[5]})
}

func (g *Graph) accept(sub string, pred string, obj object) error {
	if g == nil || g.m == nil {
		return nil
	}

	predicates, ok := g.m[sub]
	if !ok {
		o := make([]object, 0, 1)
		o = append(o, obj)
		p := make(map[string][]object)
		p[pred] = o
		g.m[sub] = p
		return nil
	}

	objects, ok := predicates[pred]
	if !ok {
		o := make([]object, 0, 1)
		o = append(o, obj)
		g.m[sub][pred] = o
		return nil
	}

	if !contains(objects, obj) {
		g.m[sub][pred] = append(g.m[sub][pred], obj)
	}

	return nil
}

func contains(objects []object, s object) bool {
	for _, object := range objects {
		if object == s {
			return true
		}
	}

	return false
}

// Bytes returns the so far consumed triples as a byte slice of
// Turle data. The triples in the byte slice are sorted first
// by subject, then by predicates, then by objects alphabetically.
func (g *Graph) Bytes() ([]byte, error) {
	if g == nil || g.m == nil {
		return nil, nil
	}

	var b []byte

	g.writePragmas(&b)

	subjects := g.sortSubjects()
	for _, subject := range subjects {
		b = append(b, []byte(fmt.Sprintf("%s ", g.sanitize(subject, "iri", false)))...)

		predicates := sortPredicates(g.m[subject])

		var predicateCounter int
		for _, predicate := range predicates {
			predicateCounter++
			objects := g.m[subject][predicate]
			sort.Slice(objects, func(i, j int) bool {
				return objects[i].item < objects[j].item
			})

			// when single predicate for a subject
			if len(predicates) == 1 {
				// write the predicate
				b = append(b, []byte(fmt.Sprintf("%s ", g.sanitize(predicate, "iri", true)))...)
				// write the predicate's objects
				g.writeObjects(&b, objects)
				continue
			}

			// when multiple predicates for subject write predicate on a new line with indentation
			b = append(b, []byte(fmt.Sprintf("\n\t%s ", g.sanitize(predicate, "iri", true)))...)

			// write the predicate's objects
			g.writeObjects(&b, objects)

			// when predicate not last, write semicolon
			if predicateCounter != len(predicates) {
				b = append(b, []byte(" ;")...)
				continue
			}
		}

		b = append(b, []byte(" .\n")...)
	}
	return b, nil
}

func (g *Graph) writeObjects(b *[]byte, objects []object) {
	for i, object := range objects {
		*b = append(*b, []byte(g.sanitizeObject(object))...)
		// when single object for predicate
		if len(objects) == 1 {
			break
		}

		// when multiple objects for predicate
		if i == (len(objects) - 1) {
			continue
		}
		*b = append(*b, []byte(", ")...)
	}
}

func (g *Graph) writePragmas(b *[]byte) {
	if g.options.Base != "" {
		*b = append(*b, []byte(fmt.Sprintf("@base <%s> .\n", g.options.Base))...)
	}

	for tag, url := range g.options.Prefixes {
		*b = append(*b, []byte(fmt.Sprintf("@prefix %s: <%s> .\n", tag, url))...)
	}
}

func (g *Graph) sortSubjects() []string {
	if g == nil || g.m == nil {
		return nil
	}
	sortedSubjects := make(sort.StringSlice, 0, len(g.m))
	for subject := range g.m {
		sortedSubjects = append(sortedSubjects, subject)
	}
	sort.Sort(sortedSubjects)
	return sortedSubjects
}

func sortPredicates(predicates map[string][]object) []string {
	sortedPredicates := make(sort.StringSlice, 0, len(predicates))
	for predicate := range predicates {
		sortedPredicates = append(sortedPredicates, predicate)
	}
	sort.Sort(sortedPredicates)
	return sortedPredicates
}
