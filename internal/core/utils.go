package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

    "github.com/nexidian/gocliselect"
	"github.com/alecthomas/chroma/lexers"
	"github.com/TyPeterson/TermJot/models"
    tm "github.com/buger/goterm"
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

// ------------- ColorText -------------
// color single words/tokens to a specific color
func ColorText(text, color string) string {
	colorCode, exists := models.ColorsMap[color]
	if !exists {
		return text
	}

	return fmt.Sprintf("\033[38;5;%dm%s\033[38;5;15m\033[48;2;57;54;70m", colorCode, text)
}

// ------------- colorTokens -------------
// color all the tokens in a code block to their specific colors
func colorTokens(text, lang string) string {
	lexer := lexers.Get(lang)
	it, err := lexer.Tokenise(nil, text)
	var finalFormattedCode string

	if err != nil {
		fmt.Println("Error tokenizing:", err)
		return ""
	}

	for _, token := range it.Tokens() {
		color := models.TokenColors[token.Type]
		finalFormattedCode += ColorText(token.Value, color)
	}

	return finalFormattedCode
}

// ------------- padLine -------------
func padLine(line string) string {
	paddingLen := 150 - len(line)
	spacesStr := strings.Repeat(" ", paddingLen)
	finalStr := fmt.Sprintf("%s%s", line, spacesStr)
	return finalStr
}

// ------------- extractCodeBlocks -------------
func extractCodeBlocks(text string) string {
	re := regexp.MustCompile("(?s)```\\s*(\\w+)?(.*?)```")

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			fmt.Println("Error extracting code block - no word after opening backticks")
			return match
		}
		lang := submatches[1]
		if lang == "" {
			fmt.Println("Error extracting code block - no language specified")
			return match
		}
		codeBlock := submatches[2]
		coloredBlock := colorTokens(codeBlock, lang)

		// add indentation to each line within code block
		lines := strings.Split(coloredBlock, "\n")
		for i, line := range lines {
			lines[i] = ColorText("      ", "white") + line
		}
		coloredBlock = strings.Join(lines[1:], "\n")

		return "\033[48;2;57;54;70m\n\n" + coloredBlock + "\033[0m"
	})

	return result
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
	text = replaceTabs(text, 4)
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


