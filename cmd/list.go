package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [[-d] [category]] | [-g]",
	Short: "List terms or categories",
	Run: func(cmd *cobra.Command, args []string) {
		// check if the -g flag is used with any other flag
		if categories && (done || len(args) > 0) {
			fmt.Println("Error: The -g flag must be used alone")
			return
		}


		if categories {
            core.ListAllCategories()
		} else {
            if len(args) > 0 {
                category := args[0]
                core.ListCategoryTerms(category, done, 111)
            } else {
                core.ListAllTerms(done)
            }
		}
	},
}
