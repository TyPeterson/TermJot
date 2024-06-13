package cmd

import (
	"fmt"
	"github.com/TyPeterson/TermJot/internal/api"
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"strings"
	"sync"
	"time"
)

const NL = "\n"

// ------------- showLoading -------------
func showLoading(done chan bool) {
	animation := []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"}

	i := 0

	// hide cursor
	fmt.Print("\033[?25l")

	defer fmt.Print("\033[?25h") // reshow cursor after function returns

	for {
		select {
		case <-done:
			fmt.Print("\033[K")
			return
		default:
			fmt.Printf("\r%s%s%s %s\t\t\t\t\t", "\033[38;5;201m", animation[i%len(animation)], "\033[0m", "Loading...")
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ------------- printWithMargins -------------
func printWithMargins(text string, margin int) {
	width, _, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	// just manually set margin to 25% of the width
	margin = int(float64(width) * 0.25)

	text = strings.TrimLeft(text, "\n")
	leftMargin := strings.Repeat(" ", margin)

	currentLineCount := 0
	// count word by word, and if currentLineCount + word.length > (width - margin), then print newline
	words := strings.Split(text, " ")
	fmt.Printf(leftMargin)

	for _, word := range words {
		if currentLineCount+len(word) > (width - (margin * 2)) {
			fmt.Printf("%s%s", NL, leftMargin)
			currentLineCount = 0
		}
		fmt.Print(word + " ")
		currentLineCount += len(word) + 1
	}

}

// ------------- displayTextWithSprite -------------
func displayTextWithSprite(text string) {
	// fmt.Println("")
	lines := strings.Split(text, "\n")

	sprite := "⚪"

	// hide cursor
	fmt.Print("\033[?25l")

	defer fmt.Print("\033[?25h") // show cursor after function returns

	for _, line := range lines {
		words := strings.Split(line, " ")
		for i, word := range words {

			fmt.Print(sprite)

			if word != "" {
				time.Sleep(20 * time.Millisecond) // (longer wait makes sprite less 'blink-y')
			}
			// move back and clear the sprite
			if i == len(words)-1 {
				fmt.Print("\033[D \033[D")
				fmt.Print(word)
			} else {
				fmt.Print("\033[D\033[D")
				fmt.Print(word + " ")
			}
		}
		fmt.Printf("\n")
	}
}

var askCmd = &cobra.Command{
	Use:   "ask [-d | -e] [-c category] term",
	Short: "Ask about term using the Gemini-1.5-Flash API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]

		api.InitializeGeminiClient()

		done := make(chan bool)
		var wg sync.WaitGroup

		definitionResult := make(chan string)
		exampleResult := make(chan string)

		// start showing loading animation
		go showLoading(done)

		// start first goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			definitionResult <- api.GenerateDefinition(term, category)
		}()

		// start second goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			exampleResult <- api.GenerateExample(term, category)
		}()

		// wait for all goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()

		r1 := core.FormatMarkdown(<-definitionResult)
		r2 := core.FormatMarkdown(<-exampleResult)

		definitionHeader := core.GenerateHeader("Description", true)
		exampleHeader := core.GenerateHeader("Example", true)

		fmt.Println("\n" + definitionHeader + "\n")
		printWithMargins(r1, 20)

		fmt.Println("\n" + exampleHeader + "\n")

		core.PrintCodeBlock(r2)

		fmt.Println("\n\n")

		// defer time.Sleep(500 * time.Millisecond) // sleep a little before exiting

	},
}
