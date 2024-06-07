package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

    "github.com/nexidian/gocliselect"
	// "github.com/alecthomas/chroma/lexers"
	"github.com/TyPeterson/TermJot/models"
    tm "github.com/buger/goterm"
    "golang.org/x/term"
)

// ------------- promptForInput -------------
func promptForInput(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ------------- promptForConfirmation -------------
func promptForConfirmation(label string) bool {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// ------------- fetchTerms -------------
func fetchTerms(category string) []models.Term {
	if category != "" {
		for _, cat := range categories {
			if cat.Name == category {
				return cat.Terms
			}
		}
	}

	return terms
}

// ------------- sortByCategory -------------
func sortByCategory(terms []models.Term) map[string][]models.Term {
    termMap := make(map[string][]models.Term)
    for _, term := range terms {
        termMap[term.Category] = append(termMap[term.Category], term)
    }
    return termMap
}


// ------------- selectTerm -------------
func selectTerm() (string, string) {
    menu := gocliselect.NewMenu("Select a term")
    numCategories := len(categories)
    if numCategories == 0 {
        fmt.Println("No categories found")
        return "", ""
    }

    // adding all terms without a category
    for _, term := range terms {
        if term.Category == "" {
            menu.AddItem(tm.Color(fmt.Sprintf("[    ]\t\t%s",  term.Name), tm.WHITE), fmt.Sprintf("%s-%s", term.Name, term.Category))
        }
    }


    for idx, category := range categories {
        for _, term := range category.Terms {
            color := 256 / numCategories * (idx + 1)
            formattedCategoryName := fmt.Sprintf("\033[38;5;%dm%s", color, category.Name)
            menu.AddItem(fmt.Sprintf("\033[38;5;15m[%s\033[38;5;15m]\033[0m\t\t%s", formattedCategoryName, term.Name), fmt.Sprintf("%s-%s", term.Name, term.Category))
        }
    }

    // -------- goterm print test --------
    // tm.Println(tm.Background(tm.Color(tm.Bold("Important header"), tm.RED), tm.WHITE))

    choice := menu.Display()
    splitChoice := strings.Split(choice, "-")

    return splitChoice[0], splitChoice[1]
}

// ------------- readKey -------------
func readKey() string {
	buf := make([]byte, 3)
	os.Stdin.Read(buf)
	return string(buf)
}

// ------------- clearScreen -------------
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ------------- filterTermsByName -------------
func filterTermsByName(terms []models.Term, name string) []models.Term {
	var filtered []models.Term
	for _, term := range terms {
		if strings.EqualFold(term.Name, name) {
			filtered = append(filtered, term)
		}
	}
	return filtered
}

// ------------- drawBoxHeader -------------
func drawBoxHeader(headerText string) {
	width, _, err := term.GetSize(0)
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	// box width is len(headerText) + 2*headerPadding
	headerPadding := 10
	boxWidth := len(headerText) + (2 * headerPadding)

	leftPadding := (width - (boxWidth + 2)) / 2

	topBorder := strings.Repeat(" ", leftPadding) + "┌" + strings.Repeat("─", boxWidth) + "┐"
	centerText := strings.Repeat(" ", leftPadding) + "│" + strings.Repeat(" ", headerPadding) + headerText + strings.Repeat(" ", headerPadding) + "│"
	botBorder := strings.Repeat(" ", leftPadding) + "└" + strings.Repeat("─", boxWidth) + "┘"

	finalHeader := topBorder + NL + centerText + NL + botBorder + NL

	fmt.Println(finalHeader)

}

// ------------- extractCodeBlocks -------------
func extractCodeBlocks(text string) string {
    re := regexp.MustCompile("(?s)```\\s*(\\w+)?(.*?)```")
	submatches := re.FindAllStringSubmatchIndex(text, -1)

	var codeBlocks []string
	for _, match := range submatches {
		codeBlockWithLang := text[match[2]:match[3]] + text[match[4]:match[5]]
		lines := strings.Split(codeBlockWithLang, "\n")
		if len(lines) < 2 {
			continue
		}
		lang := lines[0]
        // add tab to each line within code block
        for i, line := range lines {
            lines[i] = "\t" + line
        }

		codeBlock := strings.Join(lines[1:], "\n")
		coloredBlock := BackgroundColor(ColorBlockTokens(LineBreak(' ') + NL + codeBlock + LineBreak(' ') + NL, lang), 232)
		codeBlocks = append(codeBlocks, coloredBlock)
	}

	var result strings.Builder
	lastIndex := 0

	for i, match := range submatches {
		result.WriteString(text[lastIndex:match[0]])
		if i < len(codeBlocks) {
			result.WriteString(codeBlocks[i])
		}
		lastIndex = match[1]
	}
	result.WriteString(text[lastIndex:])

	return result.String()
}


// ------------- replaceTabs -------------
func replaceTabs(text string, tabWidth int) string {
	return strings.ReplaceAll(text, "\t", strings.Repeat(" ", tabWidth))
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


