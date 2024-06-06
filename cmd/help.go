package cmd

import (
	"fmt"

	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var customHelpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Displays help for a command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// No specific command provided, display brief info for all commands
			fmt.Println("Available commands:")
			for _, c := range rootCmd.Commands() {
				fmt.Printf("  %s: %s\n", c.Use, c.Short)
			}
		} else {
			// Specific command provided, display detailed help for that command
			targetCmd, _, err := rootCmd.Find(args)
			if err != nil || targetCmd == nil {
				fmt.Printf("Unknown help topic: %s\n", args[0])
				return
			}
			core.Help(targetCmd)
		}
	},
}
