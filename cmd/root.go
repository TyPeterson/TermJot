package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "termjot [OPTIONS] COMMAND [ARGS]",
	Short: "TermJot is a CLI tool for managing learning terms",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.SetHelpCommand(customHelpCmd)

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(explainCmd)
	rootCmd.AddCommand(listCmd)

	rootCmd.AddCommand(testCmd)
    rootCmd.AddCommand(defineCmd)
	InitFlags()
}
