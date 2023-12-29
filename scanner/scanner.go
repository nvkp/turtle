package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const rdfTypeIRI = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"

var (
	regexDataType = regexp.MustCompile(`(\".+\")\^\^.+`)
	regexLabel    = regexp.MustCompile(`(\".+\")@.+`)
)

// Scanner uses bufio.Scanner to parse the provided byte slice word by word.
// It keeps information about prefixes and base of the provided graph and
// the next triple to be read.
type Scanner struct {
	s        *bufio.Scanner
	t        [3]string
	prefixes map[string]string
	base     string
}

// New accepts a byte slice of the Turtle data and returns a new scanner.Scanner.
func New(data []byte) *Scanner {
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Split(scanTurtle)
	return &Scanner{
		s:        s,
		prefixes: make(map[string]string),
	}
}

// Next tries to extract a next triple from the provided data, when succesful it
// stores the new triple and returns true. If not it returns false. Another calls
// to Next would also return false.
func (s *Scanner) Next() bool {
	var index int
	var triple [3]string

	for {
		if ok := s.s.Scan(); !ok {
			return false
		}
		token := s.s.Text()

		// if bumped into a prefix form, extract and store the prefix and its value
		if token == "@prefix" {
			if ok := s.s.Scan(); !ok {
				return false
			}

			prefix := s.s.Text()
			prefix = prefix[:len(prefix)-1]

			if ok := s.s.Scan(); !ok {
				return false
			}

			value := strings.Trim(s.s.Text(), "<>")

			if !strings.HasSuffix(value, "/") {
				value = fmt.Sprintf("%s/", value)
			}

			s.prefixes[prefix] = value
			continue
		}

		// if bumped into a base form, extract and store its value
		if token == "@base" {
			if ok := s.s.Scan(); !ok {
				return false
			}

			base := strings.Trim(s.s.Text(), "<>")

			if !strings.HasSuffix(base, "/") {
				base = fmt.Sprintf("%s/", base)
			}

			s.base = base
			continue
		}

		// multiple predicates of a single subject
		if token == ";" {
			triple[0] = s.t[0] // reuse subject
			index = 1
			continue
		}

		// multiple objects of a single predicate
		if token == "," {
			triple[0] = s.t[0] // reuse subject
			triple[1] = s.t[1] // reuse predicate
			index = 2
			continue
		}

		// ignore the "end of triple" keyword
		if token == "." {
			continue
		}

		// apply the stored prefixes
		for prefix, value := range s.prefixes {
			if !strings.HasPrefix(token, fmt.Sprintf("%s:", prefix)) {
				continue
			}
			i := strings.IndexAny(token, ":")
			token = fmt.Sprintf("%s%s", value, token[i+1:])
		}

		// apply the stored base
		if strings.HasPrefix(token, "<#") {
			token = fmt.Sprintf("%s%s", s.base, token[2:])
		}

		// remove data type suffix
		if regexDataType.MatchString(token) {
			token = regexDataType.ReplaceAllString(token, `$1`)
		}

		// remove label
		if regexLabel.MatchString(token) {
			token = regexLabel.ReplaceAllString(token, `$1`)
		}

		// replace "a" keyword with rdf:type predicate
		if token == "a" {
			token = rdfTypeIRI
		}

		triple[index] = strings.Trim(token, "<>\"")
		index++

		if index > 2 {
			s.t = triple
			return true
		}
	}
}

// Triple returns the currently stored triple.
func (s *Scanner) Triple() [3]string {
	return s.t
}

func scanTurtle(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// skip leading spaces
	start := 0
	var comment bool
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])

		// a section denoted by letter # up until the new line character
		// is considered a leading space as well
		if r == '\u0023' && !comment { // #
			comment = true
			continue
		}

		if r == '\u000A' && comment { // \n
			comment = false
			continue
		}

		if !comment && !unicode.IsSpace(r) {
			break
		}
	}

	// scan until space, marking end of word
	var literal bool
	var iri bool
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		// if we bump to space character, we return the word, unless there is a literal started
		if unicode.IsSpace(r) && !literal {
			return i + width, data[start:i], nil
		}

		if (r == '\u003B' || r == '\u002C' || r == '\u002E') && !iri && !literal { // ; , .
			// if it is first character, we return it as the word
			if i == 0 {
				return i + width, data[start : i+width], nil
			}
			// otherwise we return what is before as the word
			return i, data[start:i], nil
		}

		if r == '\u0022' { // "
			literal = !literal
		}

		if r == '\u003C' || r == '\u003E' { // < >
			iri = !iri
		}
	}

	// if we're at EOF, we have a final, non-empty, non-terminated word
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// request more data.
	return start, nil, nil
}
