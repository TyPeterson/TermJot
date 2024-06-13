package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

    "github.com/nexidian/gocliselect"
	// "github.com/alecthomas/chroma/lexers"
	// "github.com/TyPeterson/TermJot/models"
    // tm "github.com/buger/goterm"
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

// ------------- selectCategory -------------
func selectCategory() string {
    menu := gocliselect.NewMenu("Select a category")
    uniqueCategories := GetUniqueCategories(false)
    for idx, category := range uniqueCategories {
        var categoryFormatted string
        if category == "" {
            categoryFormatted = TextColor("[   ]", 15)
        } else {
            color := (idx * (256/len(uniqueCategories))) + 1
            categoryFormatted = fmt.Sprintf("%s%s%s", TextColor("[", 15), TextColor(category, color), TextColor("]", 15))
        }
        menu.AddItem(categoryFormatted, category)
    }

    return menu.Display()
}


// ------------- selectTerm -------------
func selectTerm(categoryName string) string {
    menu := gocliselect.NewMenu("Select a term")
    termOptions := GetTermsInCategory(categoryName, false)

    for _, term := range termOptions {
        menu.AddItem(fmt.Sprintf("  * %s", term.Name), term.Name)
    }

    return menu.Display()
}



// ------------- indexOfString -------------
func indexOfString(slice []string, str string) int {
    for i, s := range slice {
        if strings.ToLower(s) == strings.ToLower(str) {
            return i
        }
    }

    return -1
}


// ------------- createBoxHeader -------------
func createBoxHeader(headerText string) string {
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

	finalHeader := topBorder + NL + centerText + NL + botBorder + NL

	fmt.Println(finalHeader)

    return ""
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
        // for i, line := range lines {
            // lines[i] = "    " + line

            // indent := int(float64(WIDTH) * 0.25)
            // indentSpaces := strings.Repeat(" ", indent)
            // lines[i] = indentSpaces + line

        // }

		// codeBlock := strings.Join(lines[1:], "\n")
        codeBlock := strings.Join(lines, "\n")
        // reasign first item in codeBlock to be an empty string
        codeBlock = strings.Replace(codeBlock, lang, "", 1)

		coloredBlock := ColorBlockTokens(codeBlock, lang)
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


