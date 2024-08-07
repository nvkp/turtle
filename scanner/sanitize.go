package scanner

import (
	"fmt"
	"strings"
)

const (
	dataTypeDelimiter = "^^"
	labelDelimiter    = "@"
)

func (s *Scanner) sanitize(token string) (string, string, string) {
	var label, datatype string
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

	// replace "a" keyword with rdf:type predicate
	if token == "a" {
		token = rdfTypeIRI
	}

	// trim token
	return trim(token), label, datatype
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
