package cmd

var (
    termName   string
	category   string
	done       bool
	categories bool
	define     bool
    verbose    bool
    short      bool
)

func InitFlags() {
    addCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to add")
	addCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to add the term to")
    addCmd.Flags().BoolVarP(&define, "def", "d", false, "Provide a definition for the term")


	removeCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to remove the term from")
    removeCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to remove")

	doneCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to mark the term as done in")
    doneCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to mark as done")

	askCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to explain the term in")
    askCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Provide a verbose explanation for the term")
    askCmd.Flags().BoolVarP(&short, "short", "s", false, "Provide a short explanation for the term")

	listCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to list terms")
	listCmd.Flags().BoolVarP(&done, "done", "d", false, "List only 'done' terms")
	listCmd.Flags().BoolVarP(&categories, "categories", "g", false, "List all unique categories")


}
