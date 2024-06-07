package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [-t termName] [-c category] [-d define]",
	Short: "Add a new term to the global list or to a specified category",
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Term: %s, Category: %s, Define: %t\n", termName, category, define)
        if define {
            core.HandleDefine(termName, category)
        } else {
            core.HandleAdd(termName, category)
        }

	},
}
