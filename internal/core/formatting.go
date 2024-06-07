package core


import (
"fmt"
"strings"

"golang.org/x/term"
"github.com/TyPeterson/TermJot/models"
"github.com/alecthomas/chroma/lexers"
)

const NL = "\n"
var  WIDTH int = setWidth()

func setWidth() int {
    width, _, err := term.GetSize(0)
    if err != nil {
        return 80
    }
    return width
}


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


// ----------------- generateHeader() -----------------
func GenerateHeader(headerText string) string {
    	width, _, err := term.GetSize(0)
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return ""
	}

	// box width is len(headerText) + 2*headerPadding
	headerPadding := 10
	boxWidth := len(headerText) + (2 * headerPadding)

	leftPadding := (width - (boxWidth + 2)) / 2

	topBorder := strings.Repeat(" ", leftPadding) + "┌" + strings.Repeat("─", boxWidth) + "┐"
	centerText := strings.Repeat(" ", leftPadding) + "│" + strings.Repeat(" ", headerPadding) + headerText + strings.Repeat(" ", headerPadding) + "│"
	botBorder := strings.Repeat(" ", leftPadding) + "└" + strings.Repeat("─", boxWidth) + "┘"

	finalHeader := "\n" + topBorder + "\n" + centerText + "\n" + botBorder + "\n"

    return finalHeader
}



// ------------- formatMarkdown -------------
func FormatMarkdown(text string) string {
	text = extractCodeBlocks(text)
	text = formatBold(text)
	text = formatItalic(text)
	text = formatUnderline(text)
	text = formatInlineCode(text)
	text = replaceTabs(text, 2)
	return text
}

// ------------- formatBold -------------
func formatBold(text string) string {
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	return re.ReplaceAllString(text, "\033[1m$1\033[0m")
}

// ------------- formatUnderline -------------
func formatUnderline(text string) string {
	re := regexp.MustCompile(`__(.*?)__`)
	return re.ReplaceAllString(text, "\033[4m$1\033[0m")
}

// ------------- formatItalic -------------
func formatItalic(text string) string {
	re := regexp.MustCompile(`\*(.*?)\*`)
	return re.ReplaceAllString(text, "\033[3m$1\033[0m")
}

// ------------- formatInlineCode -------------
func formatInlineCode(text string) string {
	re := regexp.MustCompile("`([^`]*)`")
	return re.ReplaceAllString(text, "\033[22m$1\033[22m")
}

// ------------- formatFaint -------------
func formatFaint(text string) string {
    return fmt.Sprintf("\033[2m%s\033[0m", text)
}

// ------------- formatWithMargins -------------
func formatWithMargins(text string, margin int) {
	width, _, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

    text = strings.TrimLeft(text, "\n")
	leftMargin := strings.Repeat(" ", margin)

	currentLineCount := 0
	// count word by word, and if currentLineCount + word.length > (width - margin), then print newline
	words := strings.Split(text, " ")
	fmt.Printf(leftMargin)

	for _, word := range words {
		if currentLineCount + len(word) > (width - (margin*2)) {
			fmt.Printf("%s%s", NL, leftMargin)
			currentLineCount = 0
		}
		fmt.Print(word + " ")
		currentLineCount += len(word) + 1
	}

}

