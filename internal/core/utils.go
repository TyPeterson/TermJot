package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
    "time"

    "github.com/nexidian/gocliselect"
	// "github.com/alecthomas/chroma/lexers"
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
        time.Sleep(50 * time.Millisecond)
    }

    fmt.Print("\n\n")
}

