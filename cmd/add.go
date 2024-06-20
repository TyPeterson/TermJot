package cmd

import (
	"errors"
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [category] [-t termName | -d define]",
	Short: "Add a new term to the global list or to a specified category",
	RunE: func(cmd *cobra.Command, args []string) error {
		if termName != "" && define {
			fmt.Println("Error: The -t and -d flags cannot be used together")
			return errors.New("Error: The -t and -d flags cannot be used together")
		}

	  if cmd.Flags().Changed("termName") && termName == "" {
            return errors.New("Error: The -t flag requires a non-empty term name")
        }

		category := ""
		if len(args) > 0 {
			category = args[0]
		}

		if define {
			core.HandleDefine(termName, category)
		} else {
			core.HandleAdd(termName, category)
		}

		return nil
	},
}
