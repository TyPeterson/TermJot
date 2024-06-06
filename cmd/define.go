package cmd

import (
    "fmt"
   "github.com/TyPeterson/TermJot/internal/core"
    "github.com/spf13/cobra"
)



var defineCmd = &cobra.Command{
    Use: "define [-c category] term definition",
    Short: "Add definition to existing term",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        term := args[0]
        definition := args[1]
        // category := cmd.Flag("category").Value.String()

        fmt.Printf("Term: %s\nDefinition: %s\n", term, definition)

        core.AddDefinition(term, definition, category)



    },

}
