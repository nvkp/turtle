package scanner

import (
	"slices"
	"unicode"
	"unicode/utf8"
)

func splitTurtle(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// skip leading spaces
	start := 0
	var comment bool
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])

		// a section denoted by letter # up until the new line character
		// is considered a leading space as well
		if r == runeNumber && !comment { // #
			comment = true
			continue
		}

		if r == runeNewLine && comment { // \n
			comment = false
			continue
		}

		if !comment && !unicode.IsSpace(r) {
			break
		}
	}

	// scan until space, marking end of word
	var literal bool
	var apostrophe bool
	var quotationMark bool
	var iri bool
	var prefixedIri bool
	var runeBuffer []rune
	var inMultiLineLiteral bool
	var backslashes int
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		// add rune to rune buffer
		runeBuffer = appendRuneToBuffer(r, runeBuffer)

		escaped := backslashes%2 == 1
		if r == runeBackslash {
			backslashes++
		} else {
			backslashes = 0
		}

		multilineLiteralEdge := bufferContainsLiterals(runeBuffer)
		// if the last characters were literals, switch the multiline
		// literal and literal state
		if multilineLiteralEdge && !escaped {
			if !inMultiLineLiteral && !literal {
				inMultiLineLiteral = true
				literal = true
			} else if inMultiLineLiteral {
				inMultiLineLiteral = false
				literal = false
			}
		}

		// if we bump to space character, we return the word, unless there is a literal started
		if unicode.IsSpace(r) && !literal {
			return i + width, data[start:i], nil
		}

		// if prefixed iri and one of the key characters and not literal and number does not follow
		// set the prefixed uri state to false
		if slices.Contains(keyCharacters, r) && !iri && !literal && prefixedIri {
			after, _ := utf8.DecodeRune(data[i+width:])

			if !unicode.IsDigit(after) {
				prefixedIri = false
			}
		}

		// if dot of a float (after it number) and not in iri and not in literal
		// return the float number
		if r == runeFullStop && !iri && !literal && !prefixedIri {
			after, afterWidth := utf8.DecodeRune(data[i+width:])

			if unicode.IsDigit(after) {
				width = width + afterWidth
				var datatype bool
				for {
					nextRune, runeWidth := utf8.DecodeRune(data[i+width:])

					if nextRune == runeCaret {
						datatype = true
					}

					if !datatype && !unicode.IsDigit(nextRune) && !slices.Contains(numberCharacters, nextRune) {
						break
					}

					if datatype && (slices.Contains(keyCharacters, nextRune) || unicode.IsSpace(nextRune)) {
						break
					}
					width = width + runeWidth
				}

				return i + width, data[start : i+width], nil
			}
		}

		if slices.Contains(keyCharacters, r) && !iri && !literal && !prefixedIri { // ; , . [
			// if it is first character, we return it as the word
			if i == 0 || start == i {
				return i + width, data[start : i+width], nil
			}
			// otherwise we return what is before as the word
			return i, data[start:i], nil
		}

		// if bumbed into quotation mark and not in apostrophe literal,
		// switch the literal and quotation mark state
		if r == runeQuotation && !apostrophe && !inMultiLineLiteral && !multilineLiteralEdge && !escaped { // "
			if !quotationMark && !literal {
				literal = true
				quotationMark = true
			} else if quotationMark {
				literal = false
				quotationMark = false
			}
		}

		// if bumbed into apostrophe and not in quotation mark literal,
		// switch the literal state and quotation mark state
		if r == runeApostrophe && !quotationMark && !inMultiLineLiteral && !multilineLiteralEdge && !escaped { // '
			if !apostrophe && !literal {
				literal = true
				apostrophe = true
			} else if apostrophe {
				literal = false
				apostrophe = false
			}
		}

		if len(runeBuffer) == 1 && !literal && !unicode.IsDigit(r) {
			prefixedIri = true
		}

		// if bumbed into the border of IRI and not in literal, switch the IRI state
		if (r == runeLessThan || r == runeGreaterThan) && !literal { // < >
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
