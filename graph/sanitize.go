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

func sanitize(str string) string {
	if len(str) == 0 {
		return str
	}

	if isIRI(str) {
		return fmt.Sprintf("<%s>", str)
	}

	if !isBlankNode(str) {
		literalEdge := `"`

		if strings.ContainsRune(str, runeNewLine) {
			if strings.ContainsRune(str, runeApostrophe) {
				literalEdge = `"""`
			} else {
				literalEdge = `'''`
			}
		}

		return fmt.Sprintf("%s%s%s", literalEdge, str, literalEdge)
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
