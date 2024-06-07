package core


import (
    "fmt"
    "strings"

    "golang.org/x/term"
    "github.com/TyPeterson/TermJot/models"
    "github.com/alecthomas/chroma/lexers"
)

const NL = "\n"

// ------------- TextColor -------------
func TextColor(text string, color int) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[38;5;15m", color, text)
}

// ------------- BackgroundColor -------------
func BackgroundColor(text string, color int) string {
	return fmt.Sprintf("\033[48;5;%dm%s\033[0m", color, text)
}

// ------------- BackgroundColorRBG -------------
func BackgroundColorRBG(text string, r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}


// ------------- ReplaceTabs -------------
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


	for _, token := range it.Tokens() {

		color := models.TokenColors[token.Type]
		coloredToken := TextColor(token.Value, models.ColorsMap[color])

		finalColoredBlock += coloredToken
	}

    // fmt.Println(BackgroundColor(LineBreak(' '), 232))
    // realFinalColoredBlock :=
//       " " +
//       NL +
        // LineBreak(' ')  +
        // finalColoredBlock

   return BackgroundColor(finalColoredBlock, 235) + NL + LineBreak(' ')

    // return  LineBreak(' ') + NL + finalColoredBlock  + NL
}
