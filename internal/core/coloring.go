package core


import (
    "fmt"
    "strings"

    "golang.org/x/term"
    "github.com/TyPeterson/TermJot/models"
    "github.com/alecthomas/chroma/lexers"
)

const NL = "\n"

// ------------- fgColor -------------
func TextColor(text string, color int) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[38;5;15m", color, text)
}

// ------------- bgColor -------------
func BackgroundColor(text string, color int) string {
	return fmt.Sprintf("\033[48;5;%dm%s\033[0m", color, text)
}

// ------------- bgColorRBG -------------
func BackgroundColorRBG(text string, r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}


// ------------- replaceTabs -------------
func ReplaceTabs(text string, tabSize int) string {
	return strings.ReplaceAll(text, "\t", strings.Repeat(" ", tabSize))
}

// ------------- LineBreak -------------
func LineBreak(char rune) string {
    width, _, err := term.GetSize(0)
    if err != nil {
        return ""
    }
    return strings.Repeat(string(char), width)
}


// ------------- ColorBlockTokens -------------
func ColorBlockTokens(text, lang string) string {
	var finalColoredBlock string

	lexer := lexers.Get(lang)
	it, err := lexer.Tokenise(nil, text)

	if err != nil {
		fmt.Println("Error tokenizing:", err)
		return ""
	}

	lineLens := []int{}
	numTokenInLine := []int{}

	curNumTokens := 0
	curLineLen := 0

	for _, token := range it.Tokens() {

		color := models.TokenColors[token.Type]
		coloredToken := TextColor(token.Value, models.ColorsMap[color])
		curLineLen += len(token.Value)
		curNumTokens++

		if token.Value == NL {
			lineLens = append(lineLens, curLineLen-1)
			numTokenInLine = append(numTokenInLine, curNumTokens)

			curLineLen = 0
			curNumTokens = 0
		}
		finalColoredBlock += coloredToken
	}


	return finalColoredBlock
}