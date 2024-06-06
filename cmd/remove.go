package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [-c category] term",
	Short: "Remove a term from the global list or from a specified category",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]
		core.RemoveTerm(term, category)
		if category != "" {
			fmt.Printf("Removed term '%s' from category '%s'\n", term, category)
		} else {
			fmt.Printf("Removed term '%s' from the global list\n", term)
		}
	},
}
