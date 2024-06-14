package cmd

import (
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)


var askCmd = &cobra.Command{
	Use:   "ask [-v | -b] [-c category] prompt",
	Short: "Ask about term using the Gemini-1.5-Flash API",
    Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        // prompt := args[0]
        // combine all args into one string
        prompt := ""
        for _, arg := range args {
            prompt += arg + " "
        }
        


        core.HandleAsk(prompt, category, verbose, brief)
	},
}
