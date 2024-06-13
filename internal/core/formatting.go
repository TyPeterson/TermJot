package core


import (
    "fmt"
    "strings"
    "regexp"

    "golang.org/x/term"
    "github.com/TyPeterson/TermJot/models"
    "github.com/alecthomas/chroma/lexers"

    "time"
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
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[39m", color, text)
}

// ------------- BackgroundColor -------------
func BackgroundColor(text string, color int) string {
	return fmt.Sprintf("\x1b[48;5;%dm%s\x1b[0m", color, text)
}

// ------------- BackgroundColorRBG -------------
func BackgroundColorRBG(text string, r, g, b int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}


// ------------- replaceTabs -------------
func replaceTabs(text string, tabSize int) string {
	return strings.ReplaceAll(text, "\t", strings.Repeat(" ", tabSize))
}

// ------------- LineBreak -------------
func LineBreak(char rune) string {
    return strings.Repeat(string(char), WIDTH)
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


    return finalColoredBlock
    // return BackgroundColor(finalColoredBlock, 235) + NL + LineBreak(' ')
}


// ----------------- generateHeader() -----------------
func GenerateHeader(headerText string) string {

	// box width is len(headerText) + 2*headerPadding
	headerPadding := 10
	boxWidth := len(headerText) + (2 * headerPadding)

	leftPadding := (WIDTH - (boxWidth + 2)) / 2

	topBorder := strings.Repeat(" ", leftPadding) + "┌" + strings.Repeat("─", boxWidth) + "┐"
	centerText := strings.Repeat(" ", leftPadding) + "│" + strings.Repeat(" ", headerPadding) + headerText + strings.Repeat(" ", headerPadding) + "│"
	botBorder := strings.Repeat(" ", leftPadding) + "└" + strings.Repeat("─", boxWidth) + "┘"

	finalHeader := "\n" + topBorder + "\n" + centerText + "\n" + botBorder

    return finalHeader
}



// ------------- formatMarkdown -------------
func FormatMarkdown(text string) string {
	text = extractCodeBlocks(text)
	// text = formatBold(text)
	// text = formatItalic(text)
	// text = formatUnderline(text)
	text = replaceTabs(text, 4)
	return text
}

// ------------- formatBold -------------
func formatBold(text string) string {
	// re := regexp.MustCompile(`\*\*(.*?)\*\*`)
    return fmt.Sprintf("\x1b[1m%s\x1b[22m", text)
}

// ------------- formatFaint -------------
func formatFaint(text string) string {
    return fmt.Sprintf("\x1b[2m%s\x1b[22m", text)
}

// ------------- formatItalic -------------
func formatItalic(text string) string {
	// re := regexp.MustCompile(`\*(.*?)\*`)
    return fmt.Sprintf("\x1b[3m%s\x1b[23m", text)
}

// ------------- formatUnderline -------------
func formatUnderline(text string) string {
    // re := regexp.MustCompile(`__(.*?)__`)
    return fmt.Sprintf("\x1b[4m%s\x1b[24m", text)
}

// ------------- formatInverted -------------
func formatInverted(text string) string {
	// re := regexp.MustCompile("`([^`]*)`")
    return fmt.Sprintf("\x1b[7m%s\x1b[27m", text)
}


// ------------- formatWithMargins -------------
func formatWithMargins(text string, margin int) {

    text = strings.TrimLeft(text, "\n")
	leftMargin := strings.Repeat(" ", margin)

	currentLineCount := 0
	// count word by word, and if currentLineCount + word.length > (width - margin), then print newline
	words := strings.Split(text, " ")
	fmt.Printf(leftMargin)

	for _, word := range words {
		if currentLineCount + len(word) > (WIDTH - (margin*2)) {
			fmt.Printf("%s%s", NL, leftMargin)
			currentLineCount = 0
		}
		fmt.Print(word + " ")
		currentLineCount += len(word) + 1
	}

}


// ------------- stripAnsiCodes -------------
func stripAnsiCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

// ------------- padLine -------------
func padLine(text string) string {
	width, _, _ := term.GetSize(0)
	padding := int(float64(width) * 0.25)
    textLen := len(stripAnsiCodes(text))

	coloredPadding := BackgroundColor(strings.Repeat(" ", padding), 0)

	textRightPadding := (width - textLen) - (padding*2)
    coloredText := BackgroundColor(text+strings.Repeat(" " , textRightPadding), 235)

	return coloredPadding + coloredText + coloredPadding
}


// ------------- printCodeBlock -------------
func PrintCodeBlock(text string) {

    for _, line := range strings.Split(text, NL) {
        fmt.Println(padLine("    " + line))
        time.Sleep(10 * time.Millisecond)
    }

}


