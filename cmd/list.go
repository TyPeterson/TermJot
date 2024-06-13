package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [[-d] [-c category]] | [-g]",
	Short: "List terms or categories",
	Run: func(cmd *cobra.Command, args []string) {
		// check if the -g flag is used with any other flag
		if categories && (done || category != "") {
			fmt.Println("Error: The -g flag cannot be used with any other flag.")
			return
		}


		if categories {
            core.ListAllCategories()
		} else {
            if category != "" {
                core.ListCategoryTerms(category, done, 111)
            } else {
                core.ListAllTerms(done)
            }
		}
	},
}
