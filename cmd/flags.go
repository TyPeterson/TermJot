package cmd

var (
    termName   string
	category   string
	done       bool
	categories bool
	define     bool
    verbose    bool
    brief      bool
)

func InitFlags() {
    addCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to add")
    addCmd.Flags().BoolVarP(&define, "def", "d", false, "Provide a definition for the term")

    removeCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to remove")

    doneCmd.Flags().StringVarP(&termName, "termName", "t", "", "Specify the term to mark as done")

	askCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to explain the term in")
    askCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Provide a verbose explanation for the term")
    askCmd.Flags().BoolVarP(&brief, "brief", "b", false, "Provide a brief explanation for the term")

	listCmd.Flags().BoolVarP(&done, "done", "d", false, "List only 'done' terms")
	listCmd.Flags().BoolVarP(&categories, "categories", "g", false, "List all unique categories")

}
