package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)


var askCmd = &cobra.Command{
	Use:   "ask [-v | -s] [-c category] prompt",
	Short: "Ask about term using the Gemini-1.5-Flash API",
    Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        prompt := args[0]

        core.HandleAsk(prompt, category, verbose, short)
	},
}
