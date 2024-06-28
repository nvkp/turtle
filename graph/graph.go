package graph

import (
	"fmt"
	"sort"
)

type object struct {
	item     string
	datatype string
	label    string
}

// Graph serves as a buffer that consumes triples one by one
// and can return a byte slice containing Turtle data of all
// triples consumed.
type Graph struct {
	m map[string]map[string][]object
}

// New returns a pointer to a new intance of graph.Graph
func New() *Graph {
	return &Graph{
		m: make(map[string]map[string][]object),
	}
}

// Accept stores a new triple to the graph.
func (g *Graph) Accept(t [3]string) error {
	return g.accept(t[0], t[1], object{item: t[2]})
}

// AcceptWithAnnotations stores a new triple with eventual label and data type
// of the object literal to the graph.
func (g *Graph) AcceptWithAnnotations(t [5]string) error {
	return g.accept(t[0], t[1], object{item: t[2], label: t[3], datatype: t[4]})
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

	subjects := sortSubjects(g)
	for _, subject := range subjects {
		b = append(b, []byte(fmt.Sprintf("%s ", sanitize(subject)))...)

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
				b = append(b, []byte(fmt.Sprintf("<%s> ", predicate))...)
				// write the predicate's objects
				writeObjects(&b, objects)
				continue
			}

			// when multiple predicates for subject write predicate on a new line with indentation
			b = append(b, []byte(fmt.Sprintf("\n\t<%s> ", predicate))...)

			// write the predicate's objects
			writeObjects(&b, objects)

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

func writeObjects(b *[]byte, objects []object) {
	for i, object := range objects {
		*b = append(*b, []byte(sanitizeObject(object))...)
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

func sortSubjects(g *Graph) []string {
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
