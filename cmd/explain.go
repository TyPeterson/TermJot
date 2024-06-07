package cmd

import (
	"fmt"
	"strings"
	"time"
    "sync"

	"github.com/TyPeterson/TermJot/internal/api"
	"github.com/TyPeterson/TermJot/internal/core"
	"github.com/spf13/cobra"
)

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

// ------------- displayTextWithSprite -------------
func displayTextWithSprite(text string) {
	fmt.Println("")
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

var explainCmd = &cobra.Command{
	Use:   "explain [-d | -e] [-c category] term",
	Short: "Explain a term using the Gemini-1.5-Flash API",
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

        definitionHeader := api.GenerateHeader("Definition")
        exampleHeader := api.GenerateHeader("Example")

        fmt.Println("\n" + definitionHeader + "\n")
        displayTextWithSprite(r1)

        fmt.Println("\n" + exampleHeader + "\n")
        displayTextWithSprite(r2)


		defer time.Sleep(500 * time.Millisecond) // sleep a little before exiting
	},
}

