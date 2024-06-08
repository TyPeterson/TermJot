package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [-c category] term",
	Short: "Mark a term as done within the global list or a specific category",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]
		core.SetTermDone(term, category)
		if category != "" {
			fmt.Printf("Marked term '%s' as done in category '%s'\n", term, category)
		} else {
			fmt.Printf("Marked term '%s' as done in the global list\n", term)
		}
	},
}
