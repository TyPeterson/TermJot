package core

import (
	"fmt"
	"strings"
	"github.com/TyPeterson/TermJot/models"
    "github.com/spf13/cobra"
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

    AddTerm(termName, definition, categoryName)
}

// ------------- HandleDefine -------------
func HandleDefine(termName, categoryName string) {
    var filteredTerms []models.Term
    var selectedTermName string
    var selectedTermCategory string


    if termName != "" && categoryName != "" {
        SetDefinition(termName, categoryName, promptForInput("Definition: "))
    } else if termName != "" && categoryName == "" {
        filteredTerms = GetTermsWithName(termName)
    } else if termName == "" && categoryName != "" {
        filteredTerms = GetTermsInCategory(categoryName)
    } else {
        filteredTerms = GetSortedByCategory()
    }

    if len(filteredTerms) == 0 {
        fmt.Println("No terms found.")
        return
    }

    selectedTermName, selectedTermCategory = selectTerm(filteredTerms)
    SetDefinition(selectedTermName, selectedTermCategory, promptForInput("Definition: "))
}

// ------------- AddTerm -------------
func AddTerm(termName, categoryName, definition string) error {

	term := models.Term {
		Name:     termName,
		Definition: definition,
		Active:   true,
		Category: categoryName,
	}
	terms = append(terms, term)

	return Save()
}


// ------------- RemoveTerm -------------
func RemoveTerm(termName, categoryName string)  {
    termOptions := GetTermsWithName(termName)

    if len(termOptions) == 0 {
        fmt.Println("No terms found.")
        return
    } else if len(termOptions) > 1 {
        termName, categoryName = selectTerm(termOptions)
    }

    for i, term := range terms {
        if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
            terms = append(terms[:i], terms[i+1:]...)
            break
        }
    }

    Save()
}


// ------------- SetDefinition -------------
func SetDefinition(termName, categoryName, definition string)  {

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
func SetTermDone(termName, categoryName string)  {

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
        if strings.ToLower(term.Name) == strings.ToLower(termName)  && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
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
func GetUniqueCategories() []string {
    uniqueCategories := make(map[string]struct{})
    for _, term := range terms {
        lowerCaseCategory := strings.ToLower(term.Category)
        if _, exists := uniqueCategories[lowerCaseCategory]; !exists {
            uniqueCategories[lowerCaseCategory] = struct{}{}
        }
    }

    categories := make([]string, 0, len(uniqueCategories))
    for category := range uniqueCategories {
        categories = append(categories, category)
    }

    return categories
}


// ------------- GetSortedByCategory -------------
func GetSortedByCategory() []models.Term {
    // return a []models.Term sorted by term.Category
    uniqueCategories := GetUniqueCategories() // get all unique categories
    sortedTerms := make([]models.Term, 0, len(terms))

    for _, category := range uniqueCategories {
        categoryTerms := GetTermsInCategory(category)
        sortedTerms = append(sortedTerms, categoryTerms...)
    }

    return sortedTerms
}

// ------------- ListAllTerms -------------
func ListAllTerms(showDone, showDoneAndActive bool) {

}

// ------------- ListCategoryTerms -------------
func ListCategoryTerms(categoryName string, showDone bool, showDoneAndActive bool) {

}


// ------------- ListAllCategories -------------
func ListAllCategories() {

}


// ------------- Help -------------
func Help(command *cobra.Command) {
	fmt.Printf("Usage: %s\n", command.Use)
	fmt.Printf("Description: %s\n", command.Short)
}
