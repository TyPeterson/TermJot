package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TyPeterson/TermJot/internal/api"
)

var storage *Storage

// ------------- HandleAdd -------------
func HandleAdd(termName, categoryName string) {

	if categoryName = filterCategoryName(categoryName); categoryName == "" {
		return
	}
	if termName = promptForInput(fmt.Sprintf("\n%s %s: ", textColor(formatBold("Term"), 14), formatFaint("[Enter to cancel]"))); termName == "" {
		return
	}

	definition := promptForInput(fmt.Sprintf("\n%s %s", textColor(formatBold("Definition"), 14), formatFaint(formatItalic("[Enter to cancel]: "))))
	addTerm(termName, categoryName, definition)
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

	definition := promptForInput(fmt.Sprintf("\n%s %s", textColor(formatBold("Definition"), 14), formatFaint(formatItalic("[Enter to cancel]: "))))
	setDefinition(termName, categoryName, definition)
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
	fmt.Println("Term marked as done")
}

// ------------- HandleAsk -------------
func HandleAsk(question, categoryName, file string, verbose, short bool) {

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

	if file != "" {
		fileName, err := findFileInSubdirectories(file)

		if err != nil {
			fmt.Printf("Error finding file: %v\n", err)
			os.Exit(1)
		}

		// read in file with os.ReadFile into fileContext
		fileContext, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		// prepend file context to prompt
		prompt = "For context, here are the contents of the " + file + " file:\n" + string(fileContext) + "\n\n" + prompt
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
func addTerm(termName, categoryName, definition string) {

	if definition == "" {
		definition = "..."
	}
	if categoryName == "" {
		categoryName = "ALL"
	}

	term := Term{
		Name:       termName,
		Definition: definition,
		Active:     true,
		Category:   strings.ToUpper(categoryName),
	}

	if err := storage.SaveData(term); err != nil {
		fmt.Printf("Error adding term: %v\n", err)
	}
}

// ------------- removeTerm -------------
func removeTerm(termName, categoryName string) {

	term, err := getTerm(termName, categoryName)
	if err != nil {
		fmt.Printf("Error removing term: %v\n", err)
	}

	if err := storage.RemoveData(term); err != nil {
		fmt.Printf("Error removing term: %v\n", err)
	}
}

// ------------- setDefinition -------------
func setDefinition(termName, categoryName, definition string) {
	term, _ := getTerm(termName, categoryName)
	term.Definition = definition

	if err := storage.UpdateData(term); err != nil {
		fmt.Printf("Error updating definition: %v\n", err)
	}
}

// ------------- setTermDone -------------
func setTermDone(termName, categoryName string) {
	term, _ := getTerm(termName, categoryName)
	term.Active = false

	if err := storage.UpdateData(term); err != nil {
		fmt.Printf("Error updating term: %v\n", err)
	}
}

// ------------- getTerm -------------
func getTerm(termName, categoryName string) (Term, error) {
	terms, err := storage.LoadAllData()
	if err != nil {
		return Term{}, err
	}

	for _, term := range terms {
		if strings.ToLower(term.Name) == strings.ToLower(termName) && strings.ToLower(term.Category) == strings.ToLower(categoryName) {
			return term, nil
		}
	}

	return Term{}, fmt.Errorf("term not found")
}

// ------------- getTermsInCategory -------------
func getTermsInCategory(categoryName string, showDone bool) []Term {
	terms, _ := storage.LoadAllData()

	var categoryTerms []Term
	for _, term := range terms {
		if strings.ToLower(term.Category) == strings.ToLower(categoryName) && term.Active != showDone {
			categoryTerms = append(categoryTerms, term)
		}
	}

	return categoryTerms
}

// ------------- getUniqueCategories -------------
func getUniqueCategories(showDone bool) []string {
	terms, _ := storage.LoadAllData()

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

// ------------- findFileInSubdirectories -------------
func findFileInSubdirectories(fileName string) (string, error) {
	// search for the file in the current directory and all subdirectories
	// return the path to the file if found, otherwise return an error
	var foundPath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, fileName) {
			foundPath = path
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("file not found")
	}

	// fmt.Println("Found file at: ", foundPath)
	return foundPath, nil
}
