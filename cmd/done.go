package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [-c category] [-t term]",
	Short: "Mark a term as done within the global list or a specific category",
	Run: func(cmd *cobra.Command, args []string) {
        core.HandleDone(termName, category)
	},
}
