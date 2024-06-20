package cmd

import (
	"errors"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [category] [-t termName]",
	Short: "Mark a term as done within the global list or a specific category",
	RunE: func(cmd *cobra.Command, args []string) error {
		termNameFlag := cmd.Flag("termName")
		if termNameFlag.Changed && termNameFlag.Value.String() == "" {
			return errors.New("Error: The -t flag requires a non-empty term name")
		}

		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		core.HandleDone(termName, category)
		return nil
	},
}
