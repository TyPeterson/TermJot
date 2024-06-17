package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [-b | -v] [-c category] prompt",
	Short: "Ask about term using the Gemini-1.5-Flash API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := ""
		for _, arg := range args {
			prompt += arg + " "
		}

		core.HandleAsk(prompt, category, verbose, brief)
	},
}
