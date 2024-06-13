package cmd

import (
	// "fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [-c category] -t [term]",
	Short: "Remove a term from the global list or from a specified category",
	Run: func(cmd *cobra.Command, args []string) {
        core.HandleRemove(termName, category)
	},
}
