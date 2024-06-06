package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [-c category] term",
	Short: "Add a new term to the global list or to a specified category",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]
		// if !strings.HasPrefix(term, `"`) || !strings.HasSuffix(term, `"`) {
		// 	fmt.Println("Error: The term must be enclosed in double quotes")
		// 	return
		// }
		fmt.Printf("Term: %s\n", term)

		// remove double quotes
		// term = term[1 : len(term)-1]

		core.AddTerm(term, category)

		if category != "" {
			fmt.Printf("Added term '%s' - [%s]\n", term, category)
		} else {
			fmt.Printf("Added term '%s' to the global list\n", term)
		}
	},
}
