package core

import (
	"github.com/alecthomas/chroma"
)

// map for ANSI color codes
var ColorsMap = map[string]int{
	"black":       0,
	"white":       15,
	"red":         1,
	"yellow":      226,
	"orange":      173,
	"blue":        27,
	"lightblue":   14,
	"lightyellow": 11,
	"lightgreen":  157,
	"mintgreen":   72,
	"green":       28,
}

// map tokenType to color
var TokenColors = map[chroma.TokenType]string{
	// KEYWORDS
	chroma.Keyword:            "blue",
	chroma.KeywordConstant:    "blue",
	chroma.KeywordDeclaration: "blue",
	chroma.KeywordNamespace:   "blue",
	chroma.KeywordPseudo:      "blue",
	chroma.KeywordReserved:    "blue",
	chroma.KeywordType:        "mintgreen",

	// NAMES
	chroma.Name:                 "lightblue",
	chroma.NameAttribute:        "lightblue",
	chroma.NameBuiltin:          "lightyellow",
	chroma.NameBuiltinPseudo:    "lightblue",
	chroma.NameClass:            "mintgreen",
	chroma.NameConstant:         "lightblue",
	chroma.NameDecorator:        "lightblue",
	chroma.NameEntity:           "lightblue",
	chroma.NameException:        "lightblue",
	chroma.NameFunction:         "lightyellow",
	chroma.NameFunctionMagic:    "lightyellow",
	chroma.NameProperty:         "lightblue",
	chroma.NameLabel:            "lightblue",
	chroma.NameNamespace:        "lightblue",
	chroma.NameOther:            "lightblue",
	chroma.NameTag:              "lightblue",
	chroma.NameVariable:         "lightblue",
	chroma.NameVariableClass:    "lightblue",
	chroma.NameVariableGlobal:   "lightblue",
	chroma.NameVariableInstance: "lightblue",
	chroma.NameVariableMagic:    "lightblue",

	// STRING LITERALS
	chroma.LiteralString:          "orange",
	chroma.LiteralStringAffix:     "orange",
	chroma.LiteralStringAtom:      "orange",
	chroma.LiteralStringBacktick:  "orange",
	chroma.LiteralStringBoolean:   "orange",
	chroma.LiteralStringChar:      "orange",
	chroma.LiteralStringDelimiter: "orange",
	chroma.LiteralStringDoc:       "orange",
	chroma.LiteralStringDouble:    "orange",
	chroma.LiteralStringEscape:    "orange",
	chroma.LiteralStringHeredoc:   "orange",
	chroma.LiteralStringInterpol:  "orange",
	chroma.LiteralStringOther:     "orange",
	chroma.LiteralStringRegex:     "orange",
	chroma.LiteralStringSingle:    "orange",
	chroma.LiteralStringSymbol:    "orange",

	// NUMERIC LITERALS
	chroma.LiteralNumber:            "lightgreen",
	chroma.LiteralNumberBin:         "lightgreen",
	chroma.LiteralNumberFloat:       "lightgreen",
	chroma.LiteralNumberHex:         "lightgreen",
	chroma.LiteralNumberInteger:     "lightgreen",
	chroma.LiteralNumberIntegerLong: "lightgreen",
	chroma.LiteralNumberOct:         "lightgreen",

	// OPERATORS
	chroma.Operator:     "white",
	chroma.OperatorWord: "white",

	// PUNCTUATION
	chroma.Punctuation: "yellow",

	// COMMENTS
	chroma.Comment:            "green",
	chroma.CommentHashbang:    "green",
	chroma.CommentMultiline:   "green",
	chroma.CommentPreproc:     "green",
	chroma.CommentPreprocFile: "green",
	chroma.CommentSingle:      "green",
	chroma.CommentSpecial:     "green",

	// WHITESPACE
	chroma.GenericDeleted:    "red",
	chroma.GenericEmph:       "red",
	chroma.GenericError:      "red",
	chroma.GenericHeading:    "red",
	chroma.GenericInserted:   "green",
	chroma.GenericOutput:     "red",
	chroma.GenericPrompt:     "red",
	chroma.GenericStrong:     "red",
	chroma.GenericSubheading: "red",
	chroma.GenericTraceback:  "red",
	chroma.GenericUnderline:  "red",

	chroma.Text: "white",
}
