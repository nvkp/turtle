package scanner

import (
	"slices"
	"unicode"
	"unicode/utf8"
)

func scanTurtle(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	var runeBuffer []rune
	var inMultiLineLiteral bool
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		// add rune to rune buffer
		runeBuffer = appendRuneToBuffer(r, runeBuffer)

		multilineLiteralEdge := bufferContainsLiterals(runeBuffer)
		escaped := escapedCharacter(runeBuffer)
		// if the last characters were literals, switch the multiline
		// literal and literal state
		if multilineLiteralEdge {
			inMultiLineLiteral = !inMultiLineLiteral
			literal = !literal
		}

		// if we bump to space character, we return the word, unless there is a literal started
		if unicode.IsSpace(r) && !literal {
			return i + width, data[start:i], nil
		}

		// if dot of a float (after it number) and not in iri and not in literal
		// return the float number
		if r == runeFullStop && !iri && !literal {
			after, afterWidth := utf8.DecodeRune(data[i+width:])

			if unicode.IsDigit(after) {
				width = width + afterWidth
				for {
					digit, digitWidth := utf8.DecodeRune(data[i+width:])
					if !unicode.IsDigit(digit) {
						break
					}
					width = width + digitWidth
				}

				return i + width, data[start : i+width], nil
			}
		}

		if slices.Contains(keyCharacters, r) && !iri && !literal { // ; , . [
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
			literal = !literal
			quotationMark = !quotationMark
		}

		// if bumbed into apostrophe and not in quotation mark literal,
		// switch the literal state and quotation mark state
		if r == runeApostrophe && !quotationMark && !inMultiLineLiteral && !multilineLiteralEdge && !escaped { // '
			literal = !literal
			apostrophe = !apostrophe
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
