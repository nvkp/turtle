package graph

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

const (
	rdfTypeIRI     = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	runeNewLine    = '\u000A' // \n
	runeApostrophe = '\u0027' // '
	runeQuotation  = '\u0022' // "

	delimQuote       = `"`
	delimApos        = `'`
	delimTripleQuote = `"""`
	delimTripleApos  = `'''`
)

func (g *Graph) sanitizeObject(obj object) string {
	item := g.sanitize(obj.item, obj.typ, false)

	if obj.label != "" {
		return fmt.Sprintf("%s@%s", item, obj.label)
	}

	if obj.datatype != "" {
		return fmt.Sprintf("%s^^%s", item, obj.datatype)
	}

	return item
}

func (g *Graph) sanitize(str string, typ string, predicate bool) string {
	if isBlankNode(str) {
		return str
	}

	if typ == "iri" || (typ == "" && isIRI(str)) {
		if str == "" {
			return str
		}

		if str == "." && g.options.Base != "" {
			return g.options.Base
		}

		if str == "a" && predicate {
			return fmt.Sprintf("<%s>", rdfTypeIRI)
		}

		for key := range g.options.Prefixes {
			if strings.HasPrefix(str, key+":") {
				return str
			}
		}

		prefix, prefixIRI := "", ""
		// find the longest match
		for key, val := range g.options.Prefixes {
			if strings.HasPrefix(str, val) && len(val) > len(prefixIRI) {
				prefix, prefixIRI = key, val
			}
		}
		if prefix != "" {
			local := strings.TrimPrefix(str, prefixIRI)
			if isValidLocalPart(local) {
				return fmt.Sprintf("%s:%s", prefix, local)
			}
		}

		if g.options.Base != "" && strings.HasPrefix(str, g.options.Base) {
			if g.options.Base == str {
				str = "."
			}

			return fmt.Sprintf("<%s>", strings.TrimPrefix(str, g.options.Base))
		}

		return fmt.Sprintf("<%s>", str)
	}

	edge := literalEdge(str)
	return fmt.Sprintf("%s%s%s", edge, str, edge)
}

func isBlankNode(str string) bool {
	return strings.HasPrefix(str, "_:")
}

func isIRI(str string) bool {
	parsedURL, err := url.Parse(str)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	for _, char := range str {
		if !isValidIRIChar(char) {
			return false
		}
	}

	return true
}

func isValidIRIChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) ||
		char == '-' || char == '.' || char == '_' || char == '~' ||
		char == ':' || char == '/' || char == '?' || char == '#' ||
		char == '[' || char == ']' || char == '@' || char == '!' ||
		char == '$' || char == '&' || char == '\'' || char == '(' ||
		char == ')' || char == '*' || char == '+' || char == ',' ||
		char == ';' || char == '=' || char == '%' ||
		unicode.Is(unicode.Han, char) || unicode.Is(unicode.Hiragana, char) ||
		unicode.Is(unicode.Katakana, char) || unicode.Is(unicode.Latin, char) ||
		unicode.Is(unicode.Arabic, char) || unicode.Is(unicode.Cyrillic, char)
}

func isValidLocalPart(local string) bool {
	if local == "" {
		return true
	}

	for i, char := range local {
		if !isValidLocalPartChar(char) {
			return false
		}

		if i == 0 && char == '.' {
			return false
		}

		if i == len(local)-1 && char == '.' {
			return false
		}
	}

	return true
}

func isValidLocalPartChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) ||
		char == '-' || char == '.' || char == '_' || char == '~' ||
		unicode.Is(unicode.Han, char) || unicode.Is(unicode.Hiragana, char) ||
		unicode.Is(unicode.Katakana, char) || unicode.Is(unicode.Latin, char) ||
		unicode.Is(unicode.Arabic, char) || unicode.Is(unicode.Cyrillic, char)
}

func literalEdge(str string) string {
	hasNL := strings.ContainsRune(str, runeNewLine)
	hasQuote := strings.ContainsRune(str, runeQuotation)
	hasApos := strings.ContainsRune(str, runeApostrophe)

	if !hasNL {
		if !hasQuote {
			return delimQuote
		}
		if !hasApos {
			return delimApos
		}
		return delimTripleQuote
	}

	if !hasQuote && !hasApos {
		return delimTripleApos
	}
	if hasApos {
		return delimTripleQuote
	}
	return delimTripleApos
}
