package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/TyPeterson/TermJot/models"
	"github.com/spf13/cobra"

	// "github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
)

var terms []models.Term
var categories []models.Category
var nextTermID int = 1
var nextCategoryID int = 1

// ------------- Init() -------------
func Init() error {
	storage, err := NewStorage()
	if err != nil {
		return err
	}
	loadedTerms, loadedCategories, err := storage.LoadData()
	if err != nil {
		return err
	}

	terms = loadedTerms
	categories = loadedCategories

	for _, term := range terms {
		if term.ID >= nextTermID {
			nextTermID = term.ID + 1
		}
	}

	for _, category := range categories {
		if category.ID >= nextCategoryID {
			nextCategoryID = category.ID + 1
		}
	}

	return nil
}

// ------------- Save -------------
func Save() error {
	storage, err := NewStorage()
	if err != nil {
		return err
	}
	return storage.SaveData(terms, categories)
}

// ------------- AddTerm -------------
func AddTerm(name, definition, categoryName string) error {
	var category *models.Category
	if categoryName != "" {
		category = FindOrCreateCategory(categoryName)
	}
	term := models.Term{
		ID:       nextTermID,
		Name:     name,
        Definition: definition,
		Active:   true,
		Category: categoryName,
	}
	nextTermID++
	terms = append(terms, term)
	if category != nil {
		category.Terms = append(category.Terms, term)
	}
	// fmt.Printf("Added term '%s'\n", name)
	return Save()
}




// ------------- AddDefinition -------------
func AddDefinition(name, definition, categoryName string) error {
    var category *models.Category
    if categoryName != "" {
        category = FindOrCreateCategory(categoryName)
    }

    // get term by name
    for i, term := range terms {
        if term.Name == name && (category == nil || term.Category == category.Name) {
            terms[i].Definition = definition
            fmt.Printf("Added definition to term '%s'\n", name)
            // term attributes:
            fmt.Printf("\n>  %s\n\n", term.Name)
            fmt.Printf("Definition: %s\n", term.Definition)
            fmt.Printf("Category: %s\n", term.Category)
            fmt.Printf("Active: %t\n", term.Active)
        }
    }
    return Save()
}


// ------------- RemoveTerm -------------
func RemoveTerm(name string, categoryName string) error {
	var category *models.Category
	if categoryName != "" {
		category = FindCategoryByName(categoryName)
		if category == nil {
			fmt.Printf("Category '%s' not found\n", categoryName)
			return nil
		}
	}
	for i, term := range terms {
		if term.Name == name && (category == nil || term.Category == category.Name) {
			terms = append(terms[:i], terms[i+1:]...)
			if category != nil {
				for j, catTerm := range category.Terms {
					if catTerm.Name == name {
						category.Terms = append(category.Terms[:j], category.Terms[j+1:]...)
						break
					}
				}
			}
			fmt.Printf("Removed term '%s'\n", name)
			return Save()
		}
	}
	fmt.Printf("Term '%s' not found\n", name)
	return nil
}

// ------------- MarkTermAsDone -------------
func MarkTermAsDone(name string, categoryName string) error {
	var category *models.Category
	if categoryName != "" {
		category = FindCategoryByName(categoryName)
		if category == nil {
			fmt.Printf("Category '%s' not found\n", categoryName)
			return nil
		}
	}
	for i, term := range terms {
		if term.Name == name && (category == nil || term.Category == category.Name) {
			terms[i].Active = false
			fmt.Printf("Marked term '%s' as done\n", name)
			return Save()
		}
	}
	fmt.Printf("Term '%s' not found\n", name)
	return nil
}

// ------------- ListTerms -------------
func ListTerms(categoryName string, showDone bool, showAll bool) {
	var category *models.Category
	if categoryName != "" {
		category = FindCategoryByName(categoryName)
		if category == nil {
			fmt.Printf("Category '%s' not found\n", categoryName)
			return
		}
	}
	for _, term := range terms {
		if (category == nil || term.Category == category.Name) && (showAll || term.Active == !showDone) {
			if term.Category != "" {
                    fmt.Printf("\n>  %s - [%s]\n\n", term.Name, term.Category)
			} else {
				fmt.Printf("\n>  %s\n\n", term.Name)
			}
            fmt.Println("Definition:", term.Definition)
		}
	}
}

// ------------- ListCategories -------------
func ListCategories() {
	for _, category := range categories {
		fmt.Println(category.Name)
	}
}

// ------------- FindOrCreateCategory -------------
func FindOrCreateCategory(name string) *models.Category {
	for i, category := range categories {
		if category.Name == name {
			return &categories[i]
		}
	}
	newCategory := models.Category{
		ID:   nextCategoryID,
		Name: name,
	}
	nextCategoryID++
	categories = append(categories, newCategory)
	return &categories[len(categories)-1]
}

// ------------- FindCategoryByName -------------
func FindCategoryByName(name string) *models.Category {
	for i, category := range categories {
		if category.Name == name {
			return &categories[i]
		}
	}
	return nil
}

// ------------- Help -------------
func Help(command *cobra.Command) {
	fmt.Printf("Usage: %s\n", command.Use)
	fmt.Printf("Description: %s\n", command.Short)
}

// ----------------- bold -----------------
func formatBold(text string) string {
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	return re.ReplaceAllString(text, "\033[1m$1\033[0m")
}

// ----------------- underline -----------------
func formatUnderline(text string) string {
	re := regexp.MustCompile(`__(.*?)__`)
	return re.ReplaceAllString(text, "\033[4m$1\033[0m")
}

// ----------------- italic -----------------
func formatItalic(text string) string {
	re := regexp.MustCompile(`\*(.*?)\*`)
	return re.ReplaceAllString(text, "\033[3m$1\033[0m")
}

// ----------------- inverse -----------------
func formatInlineCode(text string) string {
	re := regexp.MustCompile("`([^`]*)`")
	// return re.ReplaceAllString(text, "\033[38;5;15;48;5;0m[$1]\033[0m")
	return re.ReplaceAllString(text, "\033[22m$1\033[22m")
}

// ------------- colorText -------------
// color single words/tokens to a specific color
func ColorText(text, color string) string {
	colorCode, exists := models.ColorsMap[color]
	if !exists {
		return text
	}

	// return fmt.Sprintf("\033[38;5;%dm%s\033[0m", colorCode, text)
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
			fmt.Println("Error extracting code block - no word after opeing backticks ")
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
			// lines[i] = padLine(lines[i])
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
	// print original text:
	// fmt.Println("Original text:")
	// fmt.Printf("%s\n---------------------------------------------\n", text, 4)

	text = extractCodeBlocks(text)
	text = formatBold(text)
	text = formatItalic(text)
	text = formatUnderline(text)
	text = formatInlineCode(text) // + "\033[48;2;45;50;80mbingbong"

	text = replaceTabs(text, 4)

	return text
}
