package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [[-a | -d] [-c category]] | [-g]",
	Short: "List terms or categories",
	Run: func(cmd *cobra.Command, args []string) {
		// check if the -g flag is used with any other flag
		if categories && (done || all || category != "") {
			fmt.Println("Error: The -g flag cannot be used with any other flag.")
			return
		}

		// check if the -a and -d flags are used together
		if done && all {
			fmt.Println("Error: The -a and -d flags cannot be used together.")
			return
		}

		if categories {
            core.ListAllCategories()
		} else {
            if category != "" {
                core.ListCategoryTerms(category, done, all, 111)
            } else {
                core.ListAllTerms(done, all)
            }
		}
	},
}
