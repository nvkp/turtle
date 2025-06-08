package scanner

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	dataTypeDelimiter = "^^"
	labelDelimiter    = "@"
)

var numberRegex = regexp.MustCompile(`^[-0-9]+(?:\.[0-9]+)?`)

func expandPrefix(token string, value string) string {
	i := strings.Index(token, ":")
	if len(token) <= i+1 {
		return ""
	} else {
		if len(token) > i+2 && (token[i+1] == '/' || token[i+1] == '#') && value[len(value)-1] == token[i+1] {
			// if characters exist for both, trim token since we were going to do that anyway
			token = token[i+2:]
		} else if !(token[i+1] == '/' || token[i+1] == '#') && !(value[len(value)-1] == '/' || value[len(value)-1] == '#') {
			// inverse, no characters; we need to add a slash
			token = fmt.Sprintf("/%s", token[i+1:])
		} else {
			// otherwise, just trim the colon
			token = token[i+1:]
		}

		return fmt.Sprintf("<%s%s>", value, token)
	}
}

func (s *Scanner) sanitize(token string) (string, string, string, string) {
	var label, datatype string
	typ := "literal"

	// apply the stored prefixes
	for prefix, value := range s.prefixes {
		if strings.HasPrefix(token, fmt.Sprintf("%s:", prefix)) {
			token = expandPrefix(token, value)
			typ = "iri"
			break
		}
	}

	// apply the stored base
	if strings.HasPrefix(token, "<") {
		typ = "iri"
		token = trim(token)
		// short path for easy ones
		if (token == "." || token == "/") && s.base != "" {
			token = s.base
		} else {
			u, _ := url.Parse(token)
			if u == nil || u.Host == "" && s.base != "" {
				// special case for blank anchors
				if s.base[len(s.base)-1] == '#' && token[0] == '#' {
					token = s.base + token[1:]
				} else {
					b, err := url.Parse(s.base)
					if err == nil {
						if token[0] == '#' {
							// if we have # on the token side, just append the token
							token = b.String() + token
						} else {
							t := b.JoinPath(token)
							if t.String() == b.String() {
								// preserve the original form (no slash possibly, String call appends one on domains)
								token = s.base
							} else {
								token = t.String()
							}
						}
					}
				}
			}
		}
	} else if strings.HasPrefix(token, `"`) || strings.HasPrefix(token, "-") || numberRegex.MatchString(token) {
		typ = "literal"

		// extract data type suffix
		lastDataTypeIndex := lastIndex(token, dataTypeDelimiter)
		if lastDataTypeIndex != -1 {
			// Split the string into two parts
			datatype = token[lastDataTypeIndex+len(dataTypeDelimiter):]
			token = token[:lastDataTypeIndex]
		}

		// extract label suffix
		lastLabelIndex := lastIndex(token, labelDelimiter)
		if lastLabelIndex != -1 {
			// Split the string into two parts
			label = token[lastLabelIndex+len(labelDelimiter):]
			token = token[:lastLabelIndex]
		}
	} else {
		typ = "iri"

		// replace "a" keyword with rdf:type predicate
		if token == "a" {
			token = rdfTypeIRI
		}
	}

	// trim token
	return trim(token), label, datatype, typ
}

var trimmedPairs = []struct {
	left  string
	right string
}{
	{
		left:  `"""`,
		right: `"""`,
	},
	{
		left:  `'''`,
		right: `'''`,
	},
	{
		left:  "<",
		right: ">",
	},
	{
		left:  "",
		right: ">",
	},
	{
		left:  `"`,
		right: `"`,
	},
	{
		left:  `'`,
		right: `'`,
	},
}

func trim(token string) string {
	if len(token) == 0 {
		return ""
	}

	for _, pair := range trimmedPairs {
		if strings.HasPrefix(token, pair.left) && strings.HasSuffix(token, pair.right) {
			token, _ = strings.CutPrefix(token, pair.left)
			token, _ = strings.CutSuffix(token, pair.right)
			return token
		}
	}

	return token
}

var literalDelimiters = []string{
	`"""`,
	`'''`,
	`"`,
	`'`,
}

func lastIndex(token string, annotation string) int {
	for _, delimiter := range literalDelimiters {
		if !strings.HasPrefix(token, delimiter) {
			continue
		}

		lastDelimiterIndex := strings.LastIndex(token, delimiter)
		if lastDelimiterIndex == 0 {
			continue
		}

		lastAnnotationIndex := strings.LastIndex(token, annotation)

		if lastAnnotationIndex < lastDelimiterIndex {
			continue
		}

		return lastAnnotationIndex
	}
	return -1
}
