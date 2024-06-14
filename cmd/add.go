package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [category] [-t termName] [-d define]",
	Short: "Add a new term to the global list or to a specified category",
	Run: func(cmd *cobra.Command, args []string) {
		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		if define {
			core.HandleDefine(termName, category)
		} else {
			core.HandleAdd(termName, category)
		}

	},
}
