package scanner

const bufferSize = 3

func appendRuneToBuffer(r rune, buf []rune) []rune {
	buf = append(buf, r)
	if len(buf) > bufferSize {
		buf = buf[1:]
	}

	return buf
}

func containsOnlyRune(r rune, buf []rune) bool {
	for _, bufRune := range buf {
		if bufRune != r {
			return false
		}
	}

	return true
}

func bufferContainsLiterals(buf []rune) bool {
	if len(buf) != bufferSize {
		return false
	}

	return containsOnlyRune(runeQuotation, buf) || containsOnlyRune(runeApostrophe, buf) // " '
}

func escapedCharacter(buf []rune) bool {
	if len(buf) < 2 {
		return false
	}

	return buf[len(buf)-2] == runeBackslash
}
