package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexDataType = regexp.MustCompile(`(\".+\")\^\^.+`)
	regexLabel    = regexp.MustCompile(`(\".+\")@.+`)
)

func (s *Scanner) sanitize(token string) string {
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

	// trim token
	return trim(token)
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
