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

type Scanner struct {
	s        *bufio.Scanner
	t        [3]string
	prefixes map[string]string
	base     string
}

func New(data []byte) *Scanner {
	s := bufio.NewScanner(bytes.NewReader(data))
	prefixes := make(map[string]string)
	s.Split(scanTurtle)
	return &Scanner{
		s:        s,
		prefixes: prefixes,
	}
}

func (s *Scanner) Next() bool {
	var index int
	var triple [3]string

	for {
		if ok := s.s.Scan(); !ok {
			return false
		}
		token := s.s.Text()

		// TODO comment
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

		// TODO comment
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

		if token == ";" {
			triple[0] = s.t[0] // reuse subject
			index = 1
			continue
		}

		if token == "," {
			triple[0] = s.t[0] // reuse subject
			triple[1] = s.t[1] // reuse predicate
			index = 2
			continue
		}

		if token == "." {
			continue
		}

		// TODO comment
		for prefix, value := range s.prefixes {
			if !strings.HasPrefix(token, fmt.Sprintf("%s:", prefix)) {
				continue
			}
			i := strings.IndexAny(token, ":")
			token = fmt.Sprintf("%s%s", value, token[i+1:])
		}

		// TODO comment
		if strings.HasPrefix(token, "<#") {
			token = fmt.Sprintf("%s%s", s.base, token[2:])
		}

		// TODO remove data type
		if regexDataType.MatchString(token) {
			token = regexDataType.ReplaceAllString(token, `$1`)
		}

		// TODO remove data type
		if regexLabel.MatchString(token) {
			token = regexLabel.ReplaceAllString(token, `$1`)
		}

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

func (s *Scanner) Triple() [3]string {
	return s.t
}

func scanTurtle(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
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
	// Scan until space, marking end of word.
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
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
