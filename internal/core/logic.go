package core

import (
	"fmt"
	"strings"
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
        category = FindCategory(categoryName)
	}
    if category == nil {
        category = CreateCategory(categoryName)
    }


	term := models.Term {
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
    fmt.Printf("to prove it, the new term info is:\nterm: %s\ndefinition: %s\ncategory: %s\n", termName, definition, categoryName)
}

// ------------- AddDefinition -------------
func AddDefinition(name, definition, categoryName string) error {
    var category *models.Category
    if categoryName != "" {
        category = FindCategory(categoryName)
	}
    if category == nil {
        category = CreateCategory(categoryName)
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
		category = FindCategory(categoryName)
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
		category = FindCategory(categoryName)
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
    // if categoryName is provided, only list terms in that category
    if categoryName != "" {
        category := FindCategory(categoryName)

        if category == nil {
            fmt.Printf("Category '%s' not found\n", categoryName)
            return
        }

        categoryHeader := TextColor(category.Name, 111)
        fmt.Printf("\n  [%s]\n\n", categoryHeader)

        for _, term := range terms {
            if strings.ToLower(term.Category) == strings.ToLower(categoryName) {
                if showAll || term.Active || showDone {
                    fmt.Printf("  *   \x1b[1m%s:\x1b[0m\n", term.Name)
                    if term.Definition != "" {
                        fmt.Printf("\t\x1b[3m%s\x1b[0m\n\n", term.Definition)
                    } else {
                        fmt.Printf("\t\x1b[2m\x1b[3m%s\x1b[0m\n\n", "NA")
                    }
                }
            }
        }
        fmt.Println()
        return
    }

    // if no category is specified, list all terms
    numCategories := len(categories)

    // first do terms without categories
    for _, term := range terms {
        if term.Category == "" {
            termDefToDisplay := term.Definition

            if termDefToDisplay == "" {
                 termDefToDisplay = "..."
            }

            lineToDisplay := tm.Color(fmt.Sprintf("   [    ]\t\t\x1b[3m%s:\x1b[0m \x1b[2m\x1b[3m%s\x1b[0m",  term.Name, termDefToDisplay), tm.WHITE)
            fmt.Println(lineToDisplay)
        }
    }

    // now check if we have categorized terms, and if so, display them with colored category names
    if numCategories != 0 {
        for idx, category := range categories {
            for _, term := range category.Terms {
                termDefToDisplay := term.Definition

                if termDefToDisplay == "" {
                    termDefToDisplay = "..."
                }

                color := 256 / numCategories * (idx + 1)
                formattedCategoryName := fmt.Sprintf("\033[38;5;%dm%s", color, category.Name)
                fmt.Printf("\033[38;5;15m   [%s\033[38;5;15m]\033[0m\t\t\x1b[3m%s:\x1b[0m \x1b[2m\x1b[3m%s\x1b[0m\n", formattedCategoryName, term.Name, termDefToDisplay)
            }
        }
    }

}


// ------------- ListCategories -------------
func ListCategories() {
    numCategories := len(categories)
	for idx, category := range categories {
        color := 256 / (numCategories+1) * (idx + 1)
        formattedCategoryName := fmt.Sprintf("\033[38;5;%dm%s", color, category.Name)
        fmt.Printf("\033[38;5;15m   [%s\033[38;5;15m]\033[0m\n", formattedCategoryName)
	}
}

// ------------- CreateCategory -------------
func CreateCategory(name string) *models.Category {
	newCategory := models.Category{
		ID:   nextCategoryID,
		Name: name,
	}
	nextCategoryID++
	categories = append(categories, newCategory)
	return &categories[len(categories)-1]
}

// ------------- FindCategory -------------
func FindCategory(name string) *models.Category {
	for i, category := range categories {
		if strings.ToLower(category.Name) == strings.ToLower(name) {
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
