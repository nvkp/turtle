package scanner

const (
	runeNumber             = '\u0023' // #
	runeNewLine            = '\u000A' // \n
	runeSemicolon          = '\u003B' // ;
	runeComma              = '\u002C' // ,
	runeFullStop           = '\u002E' // .
	runeQuotation          = '\u0022' // "
	runeApostrophe         = '\u0027' // '
	runeLessThan           = '\u003C' // <
	runeGreaterThan        = '\u003E' // >
	runeBackslash          = '\u005C' // \
	runeLeftSquareBracket  = '\u005B' // [
	runeRightSquareBracket = '\u005D' // ]
	runeOpeningParenthesis = '\u0028' // (
	runeClosingParenthesis = '\u0029' // )
	runeUpperCaseE         = '\u0045' // E
	runeLowerCaseE         = '\u0065' // e
	runeHyphen             = '\u002D' // -
	runePlusSign           = '\u002B' // +
)

var keyCharacters = []rune{
	runeSemicolon,
	runeComma,
	runeFullStop,
	runeLeftSquareBracket,
	runeRightSquareBracket,
	runeOpeningParenthesis,
	runeClosingParenthesis,
}

var numberCharacters = []rune{
	runeUpperCaseE,
	runeLowerCaseE,
	runeHyphen,
	runePlusSign,
}
