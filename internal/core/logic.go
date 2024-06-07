package core

import (
	"fmt"
	// "strings"
	"github.com/TyPeterson/TermJot/models"
    "github.com/spf13/cobra"
    tm "github.com/buger/goterm"
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
	return Save()
}

// ------------- handleAdd -------------
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
	if categoryName != "" {
		fmt.Printf("Added term '%s' to category '%s'\n", termName, categoryName)
	} else {
		fmt.Printf("Added term '%s'\n", termName)
	}
}

// ------------- handleDefine -------------
func HandleDefine(termName, categoryName string) {
	terms := fetchTerms(categoryName)
	if len(terms) == 0 {
		fmt.Println("No terms found")
		return
	}

	if termName == "" {
		termName, categoryName = selectTerm()
        fmt.Println("termName: ", termName)
	} else {
		filteredTerms := filterTermsByName(terms, termName)
		if len(filteredTerms) == 1 {
			termName = filteredTerms[0].Name
			categoryName = filteredTerms[0].Category
		} else if len(filteredTerms) > 1 {
			termName, categoryName = selectTerm()
		} else {
			fmt.Printf("Term '%s' not found\n", termName)
			return
		}
	}

	definition := promptForInput("Definition: ")
	AddDefinition(termName, definition, categoryName)
	fmt.Printf("Added definition to term '%s'\n", termName)
}

// ------------- AddDefinition -------------
func AddDefinition(name, definition, categoryName string) error {
	var category *models.Category
	if categoryName != "" {
		category = FindOrCreateCategory(categoryName)
	}

	for i, term := range terms {
		if term.Name == name && (category == nil || term.Category == category.Name) {
			terms[i].Definition = definition
			fmt.Printf("Added definition to term '%s'\n", name)
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
    fmt.Println(tm.Background(tm.Color("IMPORTANT", tm.RED), tm.WHITE))
	fmt.Printf("Usage: %s\n", command.Use)
	fmt.Printf("Description: %s\n", command.Short)
}
