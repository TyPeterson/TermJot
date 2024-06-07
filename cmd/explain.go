package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/TyPeterson/TermJot/internal/api"
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)



func displayTextWithSprite(text string) {
	fmt.Println("")
	lines := strings.Split(text, "\n")

	sprite := "⚪"
	// sprite := []string{"⚪", "⚫"}
	// sprite := []string{".", "o", "O", "o", "."}
	// spriteIndex := 0

	// Hide the cursor
	fmt.Print("\033[?25l")

	defer fmt.Print("\033[?25h") // Ensure the cursor is shown again after the function returns

	for _, line := range lines {
		words := strings.Split(line, " ")
		for i, word := range words {

			fmt.Print(sprite)
			// fmt.Print(sprite[spriteIndex])

			if word != "" {
				time.Sleep(20 * time.Millisecond) // (longer wait makes sprite less 'blink-y')
			}
			// move back and clear the sprite
			if i == len(words)-1 {
				// Clear sprite and print the last word without trailing space
				fmt.Print("\033[D \033[D")
				fmt.Print(word)
			} else {
				fmt.Print("\033[D\033[D")
				fmt.Print(word + " ")
			}
			// spriteIndex = (spriteIndex + 1) % len(sprite)
		}
		fmt.Printf("\n")
	}
}

var explainCmd = &cobra.Command{
	Use:   "explain [-d | -e] [-c category] term",
	Short: "Explain a term using the Gemini-1.5-Flash API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]
		var explanation string
		var err error

		api.InitializeGeminiClient() // Just call the function, no assignment

		if define {
			explanation, err = api.GenerateDefinition(term, category)
		} else if example {
			explanation, err = api.GenerateExample(term, category)
		} else {
			explanation, err = api.GenerateExplanation(term, category)
		}

		if err != nil {
			fmt.Printf("Error generating content: %v\n", err)
			return
		}

		// fmt.Println(core.FormatMarkdown(explanation))		<--- boring method (booooo)

		finalExplanation := core.FormatMarkdown(explanation)
		displayTextWithSprite(finalExplanation)

		// // Text to be inserted
		// textString := "hello there"
		// another := "hi"

		// // ANSI escape code for red background
		// redBackground := "\033[41m"

		// // ANSI escape code to reset formatting
		// reset := strings.Repeat(" ", 30-len(textString)) + "\033[0m"
		// reset1 := strings.Repeat(" ", 30-len(another)) + "\033[0m"

		// // Combine the escape codes with the padded string
		// coloredString := redBackground + textString + reset
		// coloredString1 := redBackground + another + reset1

		// // Print the colored string
		// fmt.Println()
		// fmt.Println(coloredString)
		// fmt.Println(coloredString1)

		defer time.Sleep(500 * time.Millisecond) // sleep a little before exiting
	},
}
