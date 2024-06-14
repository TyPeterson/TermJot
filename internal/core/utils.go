package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
    "time"

    "github.com/nexidian/gocliselect"
    "golang.org/x/term"
)


// ------------- promptForInput -------------
func promptForInput(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ------------- selectCategory -------------
func selectCategory() string {
    fmt.Println()
    menu := gocliselect.NewMenu("Select a category")
    uniqueCategories := GetUniqueCategories(false)
    for idx, category := range uniqueCategories {
        var categoryFormatted string
        if category == "" {
            categoryFormatted = TextColor("[   ]", 15)
        } else {
            color := (idx * (256/len(uniqueCategories))) + 1
            categoryFormatted = fmt.Sprintf("%s%s%s", TextColor("[", 15), TextColor(strings.ToUpper(category), color), TextColor("]", 15))
        }
        menu.AddItem(categoryFormatted, category)
    }

    menu.AddItem(formatFaint("cancel"), "cancel selection")

    return menu.Display()
}

// ------------- selectTerm -------------
func selectTerm(categoryName string) string {
    fmt.Println()
    menu := gocliselect.NewMenu("Select a term")
    termOptions := GetTermsInCategory(categoryName, false)

    for _, term := range termOptions {
        menu.AddItem(formatBold(term.Name), term.Name)
    }

    menu.AddItem(formatFaint("cancel"), "cancel selection")

    return menu.Display()
}

// ------------- FilterCategoryName -------------
func FilterCategoryName(categoryName string) string {

    if categoryName == "" {
        categoryName = selectCategory()
        if categoryName == "cancel selection" {
            return ""
        }
    }

    if categoryName == "." {
        categoryName = GetDirectoryName()
    }

    return categoryName
}

// ------------- FilterTermName -------------
func FilterTermName(termName, categoryName string) string {
    if termName == "" {
        termName = selectTerm(categoryName)
        if termName == "cancel selection" {
            return ""
        }
    }

    return termName
}

// ------------- GetDirectoryName -------------
func GetDirectoryName() string {
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
	re := regexp.MustCompile("(?s)```\\s*(\\w+)?(.*?)```")
	submatches := re.FindAllStringSubmatchIndex(text, -1)

	var codeBlocks []string
	for _, match := range submatches {
		// Ensure the match has enough indices
		if len(match) >= 6 && match[2] >= 0 && match[3] >= 0 && match[4] >= 0 && match[5] >= 0 && match[2] <= match[3] && match[4] <= match[5] {
			codeBlockWithLang := text[match[2]:match[3]] + text[match[4]:match[5]]
			lines := strings.Split(codeBlockWithLang, "\n")

			if len(lines) < 2 {
				continue
			}

			lang := lines[0]
			codeBlock := strings.Join(lines[1:], "\n") // Join from lines[1:] to skip the language identifier

			coloredBlock := ColorBlockTokens(codeBlock, lang)
			codeBlocks = append(codeBlocks, coloredBlock)
		}
	}

	var result strings.Builder
	lastIndex := 0

	for i, match := range submatches {
		if len(match) >= 2 && match[0] >= 0 && match[1] >= 0 && match[0] <= match[1] {
			result.WriteString(text[lastIndex:match[0]])
			if i < len(codeBlocks) {
				result.WriteString(codeBlocks[i])
			}
			lastIndex = match[1]
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
			fmt.Printf("%s %s", TextColor(animation[i%len(animation)], 201), "Loading...\r")
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ------------- PrintFinalResponse -------------
func PrintFinalResponse(response string) {
    lines := strings.Split(response, NL)
    for _, line := range lines {
        fmt.Println(line)
        time.Sleep(35 * time.Millisecond)
    }

    fmt.Print("\n\n")
}

