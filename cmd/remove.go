package cmd

import (
	// "fmt"
	"errors"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [category] [-t termName]",
	Short: "Remove a term from the global list or from a specified category",
	RunE: func(cmd *cobra.Command, args []string) error {
		termNameFlag := cmd.Flag("termName")
		if termNameFlag.Changed && termNameFlag.Value.String() == "" {
			return errors.New("Error: The -t flag requires a non-empty term name")
		}

		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		core.HandleRemove(termName, category)
		return nil
	},
}
