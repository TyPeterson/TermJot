package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/term"
)

// ------------- promptForInput -------------
func promptForInput(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ------------- categoryExists -------------
func categoryExists(categoryName string) bool {
	categories := getUniqueCategories(false)

	for _, category := range categories {
		if strings.ToLower(category) == strings.ToLower(categoryName) {
			return true
		}
	}

	return false
}

// ------------- termExists -------------
func termExists(termName, categoryName string) bool {
	_, err := GetTerm(termName, categoryName)
	if err != nil {
		return false
	}

	return true
}

// ------------- selectCategory -------------
func selectCategory() string {
	fmt.Println()
	menu := Menu{Header: "Select a category:"}
	uniqueCategories := getUniqueCategories(false)

	if len(uniqueCategories) == 0 {
		uniqueCategories = append(uniqueCategories, "ALL")
	}

	for idx, category := range uniqueCategories {
		var categoryFormatted string
		if category == "" {
			categoryFormatted = textColor("[   ]", 15)
		} else {
			color := (idx * (256 / len(uniqueCategories))) + 1
			categoryFormatted = fmt.Sprintf("%s%s%s", textColor("[", 15), textColor(strings.ToUpper(category), color), textColor("]", 15))
		}
		menu.AddItem(categoryFormatted, category)
	}

	menu.AddItem(formatFaint("cancel"), "cancel selection")

	return menu.Display()
}

// ------------- selectTerm -------------
func selectTerm(categoryName string) string {
	fmt.Println()
	menu := Menu{Header: fmt.Sprintf("Select a term from %s:", categoryName)}
	termOptions := getTermsInCategory(categoryName, false)

	for _, term := range termOptions {
		menu.AddItem(formatBold(term.Name), term.Name)
	}

	menu.AddItem(formatFaint("cancel"), "cancel selection")

	return menu.Display()
}

// ------------- filterCategoryName -------------
func filterCategoryName(categoryName string) string {

	if categoryName == "" {
		categoryName = selectCategory()
		if categoryName == "cancel selection" {
			return ""
		} else {
			return categoryName
		}
	}

	if categoryName == "." {
		return getDirectoryName()
	}

	if categoryExists(categoryName) {
		return categoryName
	}

	fmt.Println("Error: Category not found")
	return ""
}

// ------------- filterTermName -------------
func filterTermName(termName, categoryName string) string {
	if termName == "" {
		termName = selectTerm(categoryName)
		if termName == "cancel selection" {
			return ""
		}
		return termName
	}

	if termExists(termName, categoryName) {
		return termName
	}

	fmt.Println("Error: Term not found")
	return ""
}

// ------------- getDirectoryName -------------
func getDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	dirSplit := strings.Split(dir, "/")
	return dirSplit[len(dirSplit)-1]
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
	// Regular expression to match code blocks with optional language identifier
	re := regexp.MustCompile("(?s)```\\s*(\\w+)?\\s*(.*?)```")
	submatches := re.FindAllStringSubmatch(text, -1)

	var codeBlocks []string
	for _, match := range submatches {
		// Ensure the match has enough indices
		if len(match) >= 3 {
			lang := strings.TrimSpace(match[1]) // The language identifier (if provided)
			codeBlock := match[2]               // The actual code block content

			coloredBlock := colorBlockTokens(codeBlock, lang)
			codeBlocks = append(codeBlocks, coloredBlock)
		}
	}

	var result strings.Builder
	lastIndex := 0

	for _, match := range submatches {
		if len(match) >= 0 {
			start := strings.Index(text[lastIndex:], match[0])
			if start >= 0 {
				result.WriteString(text[lastIndex : lastIndex+start])
				result.WriteString(codeBlocks[0])
				codeBlocks = codeBlocks[1:]
				lastIndex += start + len(match[0])
			}
		}
	}
	result.WriteString(text[lastIndex:])

	return result.String()
}

// ------------- showLoading -------------
func showLoading(done chan bool) {

	animation := []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"}

	i := 0

	// hide cursor
	fmt.Print("\033[?25l")

	defer fmt.Print("\033[?25h") // reshow cursor after function returns

	for {
		select {
		case <-done:
			fmt.Print("\033[K")
			return
		default:
			fmt.Printf("%s %s", textColor(animation[i%len(animation)], 201), "Loading...\r")
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ------------- printFinalResponse -------------
func printFinalResponse(response string) {
	lines := strings.Split(response, NL)
	for _, line := range lines {
		fmt.Println(line)
		time.Sleep(35 * time.Millisecond)
	}

	fmt.Print("\n\n")
}
