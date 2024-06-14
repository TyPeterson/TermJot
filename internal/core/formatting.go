package core

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/TyPeterson/TermJot/models"
	"github.com/alecthomas/chroma/lexers"
	"golang.org/x/term"
)

const NL = "\n"

var WIDTH int = setWidth()
var margin = int(float64(WIDTH) * 0.15)
var marginString = strings.Repeat(" ", margin)

func setWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		return 80
	}
	return width
}

// ------------- textColor -------------
func textColor(text string, color int) string {
	return fmt.Sprintf("\x1b[38;5;%dm%s\x1b[39m", color, text)
}

// ------------- backgroundColor -------------
func backgroundColor(text string, color int) string {
	return fmt.Sprintf("\x1b[48;5;%dm%s\x1b[0m", color, text)
}

// ------------- lineBreak -------------
func lineBreak(char rune) string {
	return strings.Repeat(string(char), WIDTH)
}

// ------------- colorBlockTokens -------------
func colorBlockTokens(text, lang string) string {
	var finalColoredBlock string

	lexer := lexers.Get(lang)
    if lexer == nil {
        lexer = lexers.Fallback
    }

	it, err := lexer.Tokenise(nil, text)

	if err != nil {
		fmt.Println("Error tokenizing:", err)
		return ""
	}

	for _, token := range it.Tokens() {

		color := models.TokenColors[token.Type]
		coloredToken := textColor(token.Value, models.ColorsMap[color])
		finalColoredBlock += coloredToken
	}

	return finalColoredBlock
}

// ----------------- generateHeader -----------------
func generateHeader(headerText string) string {
	// var margin int
	var headerPadding int
	var boxWidth int

	headerTextLen := len(stripAnsiCodes(headerText))

    // add 12 to boxWidth, and -6 from left padding, to create an overhang of 6 on each side
    boxWidth = (WIDTH - (2 * margin)) + 12

    headerPadding = (boxWidth - headerTextLen) / 2
    
    var offset int

    if WIDTH % 2 == 0 {
        offset = headerTextLen % 2
    } else {
        offset = (headerTextLen + 1) % 2
    }

    marginStringShortened := marginString[:len(marginString) - 6]

    topBorder := marginStringShortened + "┌" + strings.Repeat("─", boxWidth) + "┐"
    centerText := marginStringShortened + "│" + strings.Repeat(" ", headerPadding) + headerText + strings.Repeat(" ", headerPadding + offset) + "│"
    botBorder := marginStringShortened + "└" + strings.Repeat("─", boxWidth) + "┘"

    finalHeader := topBorder + NL + centerText + NL + botBorder + NL
    return finalHeader
}


// ------------- formatMarkdown -------------
func formatMarkdown(text string) string {
	text = extractCodeBlocks(text)

	// headers
	re := regexp.MustCompile(`(?m)^## (.*)$`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		re := regexp.MustCompile(`(?m)^## (.*)$`)
		submatch := re.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return formatHeader(submatch[1])
		}
		return match
	})

	// bold text
	re = regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		re := regexp.MustCompile(`\*\*(.*?)\*\*`)
		submatch := re.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return formatBold(submatch[1])
		}
		return match
	})

	// italic text
	re = regexp.MustCompile(`\*(.*?)\*`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		re := regexp.MustCompile(`\*(.*?)\*`)
		submatch := re.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return formatItalic(submatch[1])
		}
		return match
	})

	// underlined text
	re = regexp.MustCompile(`__(.*?)__`)
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		re := regexp.MustCompile(`__(.*?)__`)
		submatch := re.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return formatUnderline(submatch[1])
		}
		return match
	})

	// inlined code strings
	re = regexp.MustCompile("`([^`]*)`")
	text = re.ReplaceAllStringFunc(text, func(match string) string {
		re := regexp.MustCompile("`([^`]*)`")
		submatch := re.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return formatInverted(submatch[1])
		}
		return match
	})

	return addMargins(text)
}


// ------------- addLineMargin -------------
func addLineMargin(line string) string {
	currentLine := marginString
	currentLineCount := len(marginString)

	words := strings.Split(line, " ")
	var result []string

	for _, word := range words {
		wordLen := len(stripAnsiCodes(word))

		if currentLineCount+wordLen > (WIDTH - margin) {
			result = append(result, currentLine)
			currentLine = marginString + word
			currentLineCount = len(marginString) + wordLen + 1
		} else {
			if currentLineCount > len(marginString) {
				currentLine += " "
				currentLineCount++
			}
			currentLine += word
			currentLineCount += wordLen
		}
	}

	result = append(result, currentLine)
	return strings.Join(result, NL)
}

// ------------- addMargins -------------
func addMargins(text string) string {
	text = strings.TrimLeft(text, NL)

	lines := strings.Split(text, NL)

	var result []string

	for _, line := range lines {
		formattedLine := addLineMargin(line)
		result = append(result, formattedLine)
	}

	return strings.Join(result, NL)
}

// ------------- formatHeader -------------
func formatHeader(text string) string {
	return fmt.Sprintf("\x1b[1m\x1b[4m%s\x1b[24m\x1b[22m", text) // Bold and underline
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
	// return fmt.Sprintf("`%s`", formatItalic(formatBold(text)))
    return formatItalic(text)
}


// ------------- stripAnsiCodes -------------
func stripAnsiCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

// ------------- padLine -------------
func padLine(text string) string {
	// width, _, _ := term.GetSize(0)
	padding := int(float64(WIDTH) * 0.25)
	textLen := len(stripAnsiCodes(text))

	coloredPadding := backgroundColor(strings.Repeat(" ", padding), 0)

	textRightPadding := (WIDTH - textLen) - (padding * 2)
	coloredText := backgroundColor(text+strings.Repeat(" ", textRightPadding), 235)

	return coloredPadding + coloredText + coloredPadding
}

// ------------- formatCodeBlock -------------
func formatCodeBlock(text string) string {
	var formattedText string
	lines := strings.Split(text, NL)

	for _, line := range lines {
		formattedText += padLine("    " + line + NL)
	}

	return formattedText
}
