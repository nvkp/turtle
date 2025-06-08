package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Options changes the behavior of the scanner. It is passed to NewWithOptions.
type Options struct {
	// If set, URLs will be shortened to their relative representation to the base.
	Base string
	// If set, URLs will be shortened to their relative representation to the
	// base of the prefix on match, and the prefix applied. Resource tags ("<>")
	// will be omitted from this representation.
	Prefixes map[string]string
}

var regexBlankNode = regexp.MustCompile(`_:.+`)

// Scanner uses bufio.Scanner to parse the provided byte slice word by word.
// It keeps information about prefixes and base of the provided graph and
// the next triple to be read.
type Scanner struct {
	options          Options
	t                [][6]string
	data             []byte
	scanByteCounter  *scanByteCounter
	s                *bufio.Scanner
	base             string
	prefixes         map[string]string
	blankNodes       map[string]struct{}
	blankNodeCounter int
	curSubject       string
	curPredicate     string
	curIndex         int
	bnLists          []blankNodeList
	colls            []collection
}

type blankNodeList struct {
	start        int
	curIndex     int
	curSubject   string
	curPredicate string
	blankNode    string
}

type collection struct {
	start        int
	curIndex     int
	curSubject   string
	curPredicate string
	items        []collectionItem
}

type collectionItem struct {
	token     string
	label     string
	datatype  string
	blankNode string
	typ       string
}

// New accepts a byte slice of the Turtle data and returns a new scanner.Scanner.
func New(data []byte) *Scanner {
	return NewWithOptions(data, Options{})
}

func NewWithOptions(data []byte, options Options) *Scanner {
	counter := &scanByteCounter{}
	s := newBufioScanner(data)
	s.Split(counter.splitFunc())

	base := options.Base
	prefixes := options.Prefixes

	if prefixes == nil {
		prefixes = make(map[string]string)
	}

	return &Scanner{
		options:         options,
		data:            data,
		scanByteCounter: counter,
		s:               s,
		t:               make([][6]string, 0),
		base:            base,
		prefixes:        prefixes,
		blankNodes:      make(map[string]struct{}),
		bnLists:         make([]blankNodeList, 0),
		colls:           make([]collection, 0),
	}
}

// Next tries to extract a next triple or multiple triples from the provided
// data, when succesful it stores the new triples and returns true. If not
// it returns false. Another calls to Next would also return false.
func (s *Scanner) Next() bool {
	// shift the "pointer" of the triple slice
	if len(s.t) > 0 {
		s.t = s.t[1:]
	}

	// if there is still a triple left, return true
	if len(s.t) > 0 {
		return true
	}

	// otherwise look for next triples
	for {
		//beforeI := s.scanByteCounter.BytesRead
		if ok := s.s.Scan(); !ok {
			return false
		}

		i := s.scanByteCounter.BytesRead

		token := s.s.Text()

		// if bumped into a prefix form, extract and store the prefix and its value
		if token == "@prefix" || strings.ToLower(token) == "prefix" {
			if ok := s.s.Scan(); !ok {
				return false
			}

			prefix := s.s.Text()

			if len(prefix) == 0 {
				continue
			}

			prefix = prefix[:len(prefix)-1]

			if ok := s.s.Scan(); !ok {
				return false
			}

			value := strings.Trim(s.s.Text(), "<>")

			s.prefixes[prefix] = value
			continue
		}

		// if bumped into a base form, extract and store its value
		if token == "@base" || strings.ToLower(token) == "base" {
			if ok := s.s.Scan(); !ok {
				return false
			}

			s.base = strings.Trim(s.s.Text(), "<>")

			continue
		}

		// multiple predicates of a single subject
		if token == ";" {
			s.curIndex = 1
			continue
		}

		// multiple objects of a single predicate
		if token == "," {
			s.curIndex = 2
			continue
		}

		// ignore the "end of triple" keyword
		if token == "." {
			s.curIndex = 0
			continue
		}

		// beginning of a blank node list
		if token == "[" {
			blankNode := s.newBlankNode()
			s.bnLists = append(s.bnLists, blankNodeList{
				start:        i,
				curSubject:   s.curSubject,
				curPredicate: s.curPredicate,
				curIndex:     s.curIndex,
				blankNode:    blankNode,
			})
			s.curSubject = blankNode
			s.curIndex = 1
			continue
		}

		// ending of a blank node list
		if token == "]" {
			if len(s.bnLists) == 0 {
				continue
			}
			list := s.bnLists[len(s.bnLists)-1]
			s.bnLists = s.bnLists[:len(s.bnLists)-1]

			newData := make([]byte, 0)
			newData = append(newData, []byte(list.blankNode)...)
			newData = append(newData, s.data[i:]...)
			s.data = newData
			s.s = newBufioScanner(s.data)
			s.scanByteCounter = &scanByteCounter{}
			s.s.Split(s.scanByteCounter.splitFunc())
			s.curSubject = list.curSubject
			s.curPredicate = list.curPredicate
			s.curIndex = list.curIndex
			continue
		}

		// beginning of a collection
		if token == "(" {
			col := collection{
				start:        i,
				curIndex:     s.curIndex,
				curSubject:   s.curSubject,
				curPredicate: s.curPredicate,
				items:        make([]collectionItem, 0),
			}

			s.colls = append(s.colls, col)

			continue
		}

		if token != ")" && s.inCollection() {
			token, label, datatype, typ := s.sanitize(token)
			item := collectionItem{
				token:     token,
				label:     label,
				datatype:  datatype,
				blankNode: s.newBlankNode(),
				typ:       typ,
			}

			s.colls[len(s.colls)-1].items = append(s.colls[len(s.colls)-1].items, item)
			continue
		}

		if token == ")" {
			if len(s.colls) == 0 {
				continue
			}

			lastCollection := s.colls[len(s.colls)-1]

			s.colls = s.colls[:len(s.colls)-1]

			for i, item := range lastCollection.items {
				// rdf first
				s.t = append(s.t, [6]string{item.blankNode, rdfFirst, item.token, item.label, item.datatype, item.typ})
				// rdf rest
				rest := rdfNil
				if i < len(lastCollection.items)-1 {
					rest = lastCollection.items[i+1].blankNode
				}
				s.t = append(s.t, [6]string{item.blankNode, rdfRest, rest, "", "", "iri"})
			}

			collectionStart := rdfNilInTurtle
			if len(lastCollection.items) > 0 {
				collectionStart = lastCollection.items[0].blankNode
			}

			newData := make([]byte, 0)
			newData = append(newData, []byte(collectionStart)...)
			newData = append(newData, s.data[i:]...)
			s.data = newData
			s.s = newBufioScanner(s.data)
			s.scanByteCounter = &scanByteCounter{}
			s.s.Split(s.scanByteCounter.splitFunc())

			s.curIndex = lastCollection.curIndex
			s.curSubject = lastCollection.curSubject
			s.curPredicate = lastCollection.curPredicate

			if len(lastCollection.items) > 0 {
				return true
			}

			continue
		}

		token, label, datatype, typ := s.sanitize(token)

		// record blank node
		if regexBlankNode.MatchString(token) {
			s.blankNodes[token] = struct{}{}
		}

		// handle subject
		if s.curIndex == 0 {
			s.curSubject = token
			s.curIndex++
			continue
		}

		// handle predicate
		if s.curIndex == 1 {
			s.curPredicate = token
			s.curIndex++
			continue
		}

		// handle object
		if s.curIndex == 2 {
			s.t = append(s.t, [6]string{s.curSubject, s.curPredicate, token, label, datatype, typ})
			s.curIndex = 0
			return true
		}
	}
}

// Triple returns the next triple
func (s *Scanner) Triple() [3]string {
	if len(s.t) == 0 {
		return [3]string{}
	}
	return [3]string{s.t[0][0], s.t[0][1], s.t[0][2]}
}

// TripleWithAnnotations returns the next triple with label and datatype
func (s *Scanner) TripleWithAnnotations() [6]string {
	if len(s.t) == 0 {
		return [6]string{}
	}
	return s.t[0]
}

// newBlankNode emits a new blank node based on what is the
// blank node ID counter and what blank node have already
// been recorded in the dataset to avoid collisions
func (s *Scanner) newBlankNode() string {
	for {
		blankNode := fmt.Sprintf("_:b%d", s.blankNodeCounter)
		s.blankNodeCounter = s.blankNodeCounter + 1
		if _, ok := s.blankNodes[blankNode]; ok {
			continue
		}

		s.blankNodes[blankNode] = struct{}{}
		return blankNode
	}
}

func (s *Scanner) inCollection() bool {
	if len(s.colls) == 0 {
		return false
	}

	if len(s.bnLists) == 0 {
		return true
	}

	return s.colls[len(s.colls)-1].start > s.bnLists[len(s.bnLists)-1].start
}

func newBufioScanner(data []byte) *bufio.Scanner {
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Buffer(data, 0)
	return s
}
