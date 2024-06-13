package core

import (
	"fmt"
	"strings"
	"github.com/TyPeterson/TermJot/internal/api"
	"github.com/TyPeterson/TermJot/models"
	"github.com/spf13/cobra"
)

var terms []models.Term

// ------------- HandleAdd -------------
func HandleAdd(termName, categoryName string) {
	if termName == "" {
		termName = promptForInput("Term: ")
	}

	if categoryName == "" {
        categoryName = promptForInput(fmt.Sprintf("%s %s", formatBold("Category"), formatFaint(formatItalic("[Enter to skip]: "))))
	}

	var definition string
definition = promptForInput(fmt.Sprintf("%s %s", formatBold("Definition"), formatFaint(formatItalic("[Enter to skip]: "))))

	AddTerm(termName, categoryName, definition)
}

// ------------- HandleDefine -------------
func HandleDefine(termName, categoryName string) {

	if categoryName == "" {
		categoryName = selectCategory()
	}

	if termName == "" {
		termName = selectTerm(categoryName)
	}

	SetDefinition(termName, categoryName, promptForInput(formatBold("Definition: ")))
}

// ------------- HandleAdd -------------
func HandleRemove(termName, categoryName string) {

	if categoryName == "" {
		categoryName = selectCategory()
	}
	if termName == "" {
		termName = selectTerm(categoryName)
	}
	fmt.Println("Removing term:", termName, "from category:", categoryName)
	RemoveTerm(termName, categoryName)
}

// ------------- HandleDone -------------
func HandleDone(termName, categoryName string) {

	if categoryName == "" {
		categoryName = selectCategory()
	}
	if termName == "" {
		termName = selectTerm(categoryName)
	}

	SetTermDone(termName, categoryName)
}

// ------------- HandleAsk -------------
func HandleAsk(question string, categoryName string, verbose, short bool) {

	prompt := question
	var responseType string

	var geminiResponse string

	if categoryName != "" {
		prompt = fmt.Sprintf("%s in context of: %s", question, categoryName)
	}

	if verbose {
		responseType = "verbose"
	} else if short {
		responseType = "short"
	} else {
		responseType = "default"
	}

	done := make(chan bool)
	go showLoading(done)

	geminiResponse = api.GetResponse(prompt, responseType)
	done <- true

	formattedResult := FormatMarkdown(geminiResponse)

	responseHeader := GenerateHeader(formatBold("J O T"), true)
	fmt.Println("\n" + responseHeader + "\n")

	PrintFinalResponse(formattedResult)

}

// ------------- AddTerm -------------
func AddTerm(termName, categoryName, definition string) error {

	if definition == "" {
		definition = "..."
	}
	if categoryName == "" {
		categoryName = "_none_"
	}

	term := models.Term{
		Name:       termName,
		Definition: definition,
		Active:     true,
		Category:   categoryName,
	}
	terms = append(terms, term)

	return Save()
}

// ------------- RemoveTerm -------------
func RemoveTerm(termName, categoryName string) {
	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			terms = append(terms[:i], terms[i+1:]...)
			break
		}
	}

	Save()
}

// ------------- SetDefinition -------------
func SetDefinition(termName, categoryName, definition string) {

	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			term.Definition = definition
			terms[i] = term
			break
		}
	}

	Save()
}

// ------------- SetTermDone -------------
func SetTermDone(termName, categoryName string) {

	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			term.Active = false
			terms[i] = term
			break
		}
	}

	Save()
}

// ------------- GetTerm -------------
func GetTerm(termName, categoryName string) models.Term {
	for _, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			return term
		}
	}

	return models.Term{}
}

// ------------- GetTermsInCategory -------------
func GetTermsInCategory(categoryName string, showDone bool) []models.Term {
	var categoryTerms []models.Term
	for _, term := range terms {
		if strings.ToLower(term.Category) == strings.ToLower(categoryName) && term.Active != showDone {
			categoryTerms = append(categoryTerms, term)
		}
	}

	return categoryTerms
}

// ------------- GetTermsWithName -------------
func GetTermsWithName(termName string) []models.Term {
	var namedTerms []models.Term
	for _, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) {
			namedTerms = append(namedTerms, term)
		}
	}

	return namedTerms
}

// ------------- GetUniqueCategories -------------
func GetUniqueCategories(showDone bool) []string {
	uniqueCategories := make(map[string]struct{})
	for _, term := range terms {
		if term.Active != showDone { // Consider terms based on the opposite of showDone
			lowerCaseCategory := strings.ToLower(term.Category)
			if _, exists := uniqueCategories[lowerCaseCategory]; !exists {
				uniqueCategories[lowerCaseCategory] = struct{}{}
			}
		}
	}

	categories := make([]string, 0, len(uniqueCategories))
	for category := range uniqueCategories {
		categories = append(categories, category)
	}

	return categories
}


// ------------- ListCategoryTerms -------------
func ListCategoryTerms(categoryName string, showDone bool, color int) {

	categoryTerms := GetTermsInCategory(categoryName, showDone)
	formattedHeader := GenerateHeader(TextColor(formatBold(strings.ToUpper(categoryName)), color), false)
	headerPrinted := false

	for _, term := range categoryTerms {
		if term.Active == showDone {
			continue
		}

		if !headerPrinted {
			fmt.Printf("\n%s\n", formattedHeader)
			headerPrinted = true
		}
        box := "☐" 
        if showDone {
            box = "☑"
        }
		termFormatted := fmt.Sprintf("   %s   %s: %s\n", box, formatBold(term.Name), formatItalic(formatFaint(term.Definition)))
		fmt.Println(strings.ReplaceAll(AddLineMargin(termFormatted), fmt.Sprintf("%s%s", NL, marginString), fmt.Sprintf("%s%s%s", NL, marginString, strings.Repeat(" ", 7))))
	}

}

// ------------- ListAllTerms -------------
func ListAllTerms(showDone bool) {
	// for each category, call ListCategoryTerms
	uniqueCategories := GetUniqueCategories(showDone)
	for idx, category := range uniqueCategories {
		color := (idx * (256 / len(uniqueCategories))) + 1
		ListCategoryTerms(category, showDone, color)
	}

}

// ------------- ListAllCategories -------------
func ListAllCategories() {
	// get all unique categories and print them out
	uniqueCategories := GetUniqueCategories(false)

	fmt.Printf("\n    %s\n\n", formatUnderline(formatBold("Categories:")))

	for idx, category := range uniqueCategories {
		color := (idx * (256 / len(uniqueCategories))) + 1
		categoryFormatted := fmt.Sprintf(" *  %s%s%s\n", TextColor("[", 15), TextColor(category, color), TextColor("]", 15))
		if category != "" {
			fmt.Println(categoryFormatted)
		}
	}

}

// ------------- Help -------------
func Help(command *cobra.Command) {
	fmt.Printf("Usage: %s\n", command.Use)
	fmt.Printf("Description: %s\n", command.Short)
}
