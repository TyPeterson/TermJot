package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [category] -t [term]",
	Short: "Remove a term from the global list or from a specified category",
	Run: func(cmd *cobra.Command, args []string) {
        category := ""
        if len(args) > 0 {
            category = args[0]
        }
        core.HandleRemove(termName, category)
	},
}
