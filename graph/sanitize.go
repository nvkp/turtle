package graph

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

const (
	runeNewLine    = '\u000A' // \n
	runeApostrophe = '\u0027' // '
	runeQuotation  = '\u0022' // "
)

func (g *Graph) sanitizeObject(obj object) string {
	item := g.sanitize(obj.item, false)

	if obj.label != "" {
		return fmt.Sprintf("%s@%s", item, obj.label)
	}

	if obj.datatype != "" {
		return fmt.Sprintf("%s^^%s", item, obj.datatype)
	}

	return item
}

func (g *Graph) sanitize(str string, predicate bool) string {
	if len(str) == 0 {
		return str
	}

	if isIRI(str) {
		if g.options.ResolveURLs {
			for key, prefix := range g.options.Prefixes {
				if strings.HasPrefix(str, prefix) {
					return fmt.Sprintf("%s:%s", key, strings.TrimPrefix(str, prefix))
				}
			}

			if g.options.Base != "" && strings.HasPrefix(str, g.options.Base) {
				if g.options.Base == str {
					str = "."
				}

				return fmt.Sprintf("<%s>", strings.TrimPrefix(str, g.options.Base))
			}
		}

		return fmt.Sprintf("<%s>", str)
	}

	if isBlankNode(str) {
		return str
	}

	if !isBlankNode(str) && !predicate {
		edge := literalEdge(str)

		return fmt.Sprintf("%s%s%s", edge, str, edge)
	}

	return str
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

// TODO consts

func literalEdge(str string) string {
	if !strings.ContainsRune(str, runeNewLine) {
		if !strings.ContainsRune(str, runeQuotation) {
			return `"`
		}

		return `'`
	}

	if strings.ContainsRune(str, runeApostrophe) {
		return `"""`
	}

	return `'''`
}
