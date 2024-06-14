package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [category] [-t term]",
	Short: "Mark a term as done within the global list or a specific category",
	Run: func(cmd *cobra.Command, args []string) {
		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		core.HandleDone(termName, category)
	},
}
