package graph

import (
	"fmt"
	"sort"
)

type Graph map[string]map[string][]string

func (g Graph) Accept(t [3]string) error {
	if g == nil {
		return nil
	}

	predicates, ok := g[t[0]]
	if !ok {
		o := make([]string, 0, 1)
		o = append(o, t[2])
		p := make(map[string][]string)
		p[t[1]] = o
		g[t[0]] = p
		return nil
	}

	objects, ok := predicates[t[1]]
	if !ok {
		o := make([]string, 0, 1)
		o = append(o, t[2])
		g[t[0]][t[1]] = o
		return nil
	}

	if !contains(objects, t[2]) {
		g[t[0]][t[1]] = append(g[t[0]][t[1]], t[2])
	}

	return nil
}

func contains(objects []string, s string) bool {
	for _, object := range objects {
		if object == s {
			return true
		}
	}

	return false
}

func (g Graph) Bytes() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	var b []byte

	subjects := g.sortSubjects()
	for _, subject := range subjects {
		b = append(b, []byte(fmt.Sprintf("<%s> ", subject))...)

		predicates := sortPredicates(g[subject])

		var predicateCounter int
		for _, predicate := range predicates {
			predicateCounter++
			objects := g[subject][predicate]
			sort.Slice(objects, func(i, j int) bool {
				return objects[i] < objects[j]
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

func writeObjects(b *[]byte, objects []string) {
	for i, object := range objects {
		// TODO literals are not printed with <>
		*b = append(*b, []byte(fmt.Sprintf("<%s>", object))...)
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

func (g Graph) sortSubjects() []string {
	sortedSubjects := make(sort.StringSlice, 0, len(g))
	for subject := range g {
		sortedSubjects = append(sortedSubjects, subject)
	}
	sort.Sort(sortedSubjects)
	return sortedSubjects
}

func sortPredicates(predicates map[string][]string) []string {
	sortedPredicates := make(sort.StringSlice, 0, len(predicates))
	for predicate := range predicates {
		sortedPredicates = append(sortedPredicates, predicate)
	}
	sort.Sort(sortedPredicates)
	return sortedPredicates
}
