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

    if categoryName = filterCategoryName(categoryName); categoryName == "" { 
        return
    }
    if termName = promptForInput(fmt.Sprintf("\n%s %s: ", textColor(formatBold("Term"), 14), formatFaint("[Enter to cancel]"))); termName == "" {
        return
    }

    addTerm(termName, categoryName, promptForInput(fmt.Sprintf("\n%s %s", textColor(formatBold("Definition"), 14), formatFaint(formatItalic("[Enter to cancel]: ")))))
    fmt.Println("\n\nTerm added successfully")
}

// ------------- HandleDefine -------------
func HandleDefine(termName, categoryName string) {

    if categoryName = filterCategoryName(categoryName); categoryName == "" {
        return
    }
    if termName = filterTermName(termName, categoryName); termName == "" {
        return
    }

    setDefinition(termName, categoryName, promptForInput(fmt.Sprintf("\n%s %s", textColor(formatBold("Definition"), 14), formatFaint(formatItalic("[Enter to cancel]: ")))))
    fmt.Println("\n\nDefinition update successful")
}

// ------------- HandleRemove -------------
func HandleRemove(termName, categoryName string) {

    if categoryName = filterCategoryName(categoryName); categoryName == "" {
        return
    }
    if termName = filterTermName(termName, categoryName); termName == "" {
        return
    }

	removeTerm(termName, categoryName)
    fmt.Println("\n\nTerm removed successfully")
}

// ------------- HandleDone -------------
func HandleDone(termName, categoryName string) {

    if categoryName = filterCategoryName(categoryName); categoryName == "" {
        return
    }
    if termName = filterTermName(termName, categoryName); termName == "" {
        return
    }

	setTermDone(termName, categoryName)
    fmt.Println("\n\nTerm marked as done")
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
		responseType = "brief"
	} else {
		responseType = "default"
	}

	done := make(chan bool)
	go showLoading(done)

	geminiResponse = api.GetResponse(prompt, responseType)
	done <- true

	formattedResult := formatMarkdown(geminiResponse)

	responseHeader := generateHeader(formatBold("J O T"))
	fmt.Println("\n" + responseHeader + "\n")

	printFinalResponse(formattedResult)
}

// ------------- addTerm -------------
func addTerm(termName, categoryName, definition string) error {

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

// ------------- removeTerm -------------
func removeTerm(termName, categoryName string) {
	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			terms = append(terms[:i], terms[i+1:]...)
			break
		}
	}

	Save()
}

// ------------- setDefinition -------------
func setDefinition(termName, categoryName, definition string) {

	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			term.Definition = definition
			terms[i] = term
			break
		}
	}

	Save()
}

// ------------- setTermDone -------------
func setTermDone(termName, categoryName string) {

	for i, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			term.Active = false
			terms[i] = term
			break
		}
	}

	Save()
}

// ------------- getTermsInCategory -------------
func getTermsInCategory(categoryName string, showDone bool) []models.Term {
	var categoryTerms []models.Term
	for _, term := range terms {
		if strings.ToLower(term.Category) == strings.ToLower(categoryName) && term.Active != showDone {
			categoryTerms = append(categoryTerms, term)
		}
	}

	return categoryTerms
}

// ------------- getUniqueCategories -------------
func getUniqueCategories(showDone bool) []string {
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
    
    if categoryName == "." {
        categoryName = getDirectoryName()
    }

	categoryTerms := getTermsInCategory(categoryName, showDone)
	formattedHeader := generateHeader(textColor(formatBold(strings.ToUpper(categoryName)), color))
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

		termFormatted := fmt.Sprintf("%s  %s: %s\n", box, formatBold(term.Name), formatItalic(formatFaint(term.Definition)))

		termFormattedWithMargins := strings.ReplaceAll(addLineMargin(termFormatted), fmt.Sprintf("%s%s", NL, marginString), fmt.Sprintf("%s%s   ", NL, marginString))
		fmt.Println(termFormattedWithMargins)
	}

}

// ------------- ListAllTerms -------------
func ListAllTerms(showDone bool) {
	// for each category, call ListCategoryTerms
	uniqueCategories := getUniqueCategories(showDone)
	for idx, category := range uniqueCategories {
		color := (idx * (256 / len(uniqueCategories))) + 1
		ListCategoryTerms(category, showDone, color)
	}

}

// ------------- ListAllCategories -------------
func ListAllCategories() {
	// get all unique categories and print them out
	uniqueCategories := getUniqueCategories(false)

	fmt.Printf("\n    %s\n\n", formatUnderline(formatBold("Categories:")))

	for idx, category := range uniqueCategories {
		color := (idx * (256 / len(uniqueCategories))) + 1
		categoryFormatted := fmt.Sprintf(" *  %s%s%s\n", textColor("[", 15), textColor(category, color), textColor("]", 15))

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
