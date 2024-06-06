package cmd

var (
	category   string
	done       bool
	all        bool
	categories bool
	define     bool
	example    bool
)

func InitFlags() {
	addCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to add the term to")
    // TODO: Add flag for '-d' or '--define'

	removeCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to remove the term from")

	doneCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to mark the term as done in")

	explainCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to explain the term in")
	explainCmd.Flags().BoolVarP(&define, "define", "d", false, "Provide a definition for the term")
	explainCmd.Flags().BoolVarP(&example, "example", "e", false, "Provide an example for the term")

	listCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to list terms")
	listCmd.Flags().BoolVarP(&done, "done", "d", false, "List only 'done' terms")
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List both 'active' and 'done' terms")
	listCmd.Flags().BoolVarP(&categories, "categories", "g", false, "List all unique categories")

    // add define flags (-c)
    defineCmd.Flags().StringVarP(&category, "category", "c", "", "Specify a category to define the term in")
    defineCmd.Flags().StringVarP(&category, "define", "d", "", "Provide a definition for the term")

}
