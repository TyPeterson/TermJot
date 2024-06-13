package core

import (
	"fmt"
	"github.com/TyPeterson/TermJot/models"
	"github.com/spf13/cobra"
	"strings"
)

var terms []models.Term

// ------------- HandleAdd -------------
func HandleAdd(termName, categoryName string) {
	if termName == "" {
		termName = promptForInput("Term: ")
	}

	if categoryName == "" && promptForConfirmation("Add to a category? (y/n): ") {
		categoryName = promptForInput("Category: ")
	}

	var definition string
	if promptForConfirmation("Add a definition? (y/n): ") {
		definition = promptForInput("Definition: ")
	}

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

	SetDefinition(termName, categoryName, promptForInput("Definition: "))
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

// ------------- AddTerm -------------
func AddTerm(termName, categoryName, definition string) error {

	if definition == "" {
		definition = "..."
	}
	if categoryName == "" {
		categoryName = "_all_"
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
func GetTermsInCategory(categoryName string) []models.Term {
	var categoryTerms []models.Term
	for _, term := range terms {
		if strings.ToLower(term.Category) == strings.ToLower(categoryName) {
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
// func GetUniqueCategories(showDone bool) []string {
//     fmt.Println("GetUniqueCategories called with showDone:", showDone)
//
//     // for _, term := range terms {
//     //     if term.Active != showDone {
//     //         fmt.Println(term.Name, term.Active)
//     //     }
//     // }
//
// 	uniqueCategories := make(map[string]struct{})
// 	for _, term := range terms {
// 		if term.Active { // Only consider terms specified by showDone
// 			lowerCaseCategory := strings.ToLower(term.Category)
// 			if _, exists := uniqueCategories[lowerCaseCategory]; !exists {
// 				uniqueCategories[lowerCaseCategory] = struct{}{}
// 			}
// 		}
// 	}
//
// 	categories := make([]string, 0, len(uniqueCategories))
// 	for category := range uniqueCategories {
// 		categories = append(categories, category)
// 	}
//
// 	return categories
// }
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

    for _, category := range categories {
        fmt.Println(category)
    }
	return categories
}


// ------------- GetSortedByCategory -------------
func GetSortedByCategory() []models.Term {

	uniqueCategories := GetUniqueCategories(false) // get all unique categories
	sortedTerms := make([]models.Term, 0, len(terms))

	for _, category := range uniqueCategories {
		categoryTerms := GetTermsInCategory(category)
		sortedTerms = append(sortedTerms, categoryTerms...)
	}

	return sortedTerms
}

// ------------- ListCategoryTerms -------------
func ListCategoryTerms(categoryName string, showDone bool, color int) {

    categoryTerms := GetTermsInCategory(categoryName)
    formattedHeader := GenerateHeader(TextColor(formatBold(strings.ToUpper(categoryName)), color), false)
    headerPrinted := false
    
    for _, term := range categoryTerms {
        if term.Active == showDone {
            continue
        }
        // fmt.Println("term:", term.Name, "active:", term.Active, "showDone:", showDone)
        if !headerPrinted {
            fmt.Println(formattedHeader)
            fmt.Println()
            headerPrinted = true
        }
        termFormatted := fmt.Sprintf("  * %s: %s\n", formatBold(term.Name), formatItalic(formatFaint(term.Definition)))
        fmt.Println(termFormatted)
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
