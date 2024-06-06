package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"os"
	// "github.com/TyPeterson/TermJot/internal/core"
	// "github.com/TyPeterson/TermJot/internal/core"
	"github.com/TyPeterson/TermJot/models"
	"github.com/google/generative-ai-go/genai"

	// "github.com/joho/godotenv"
	"golang.org/x/term"
	"google.golang.org/api/option"
)

var client *genai.Client

// ----------------- InitializeGeminiClient() -----------------
func InitializeGeminiClient() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not found in environment variables")
	}

	ctx := context.Background()
	var clientErr error
	client, clientErr = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if clientErr != nil {
		log.Fatalf("Error creating Gemini client: %v", clientErr)
	}
}

// ----------------- GenerateExplanation() -----------------
func GenerateExplanation(term, category string) (string, error) {
	// return generateContent(term, category, "explain")
	return generateContent(term, category, "both")
}

// ----------------- GenerateDefinition() -----------------
func GenerateDefinition(term, category string) (string, error) {
	return generateContent(term, category, "define")
}

// ----------------- GenerateExample() -----------------
func GenerateExample(term, category string) (string, error) {
	return generateContent(term, category, "example")
}

func constructString(char byte, spacesBeforeCount, spacesAfterCount int) string {
	spacesBeforeStr := strings.Repeat(" ", spacesBeforeCount)
	spacesAfterStr := strings.Repeat(" ", spacesAfterCount)
	return spacesBeforeStr + string(char) + spacesAfterStr
}

// ----------------- showLoading() -----------------
func showLoading(done chan bool) {

	// totalSpaces := 16
	// const char = '°'
	// animation := make([]string, (totalSpaces*2)+1)
	// animation1 := make([]string, (totalSpaces*2)+1)
	// animation2 := make([]string, (totalSpaces*2)+1)

	// // construct animation1 (left side hitter)
	// animation1[0] = "]"
	// animation1[1] = "/"
	// animation1[2] = "-"
	// animation1[3] = "/"

	// for i := 4; i < (totalSpaces * 2); i++ {
	// 	animation1[i] = "|"
	// }
	// animation1[totalSpaces*2] = "\\"

	// // construct animation2 (right side hitter)
	// for i := 0; i < totalSpaces-1; i++ {
	// 	animation2[i] = "|"
	// }
	// animation2[totalSpaces-1] = "/"
	// animation2[totalSpaces] = "["
	// animation2[totalSpaces+1] = "\\"
	// animation2[totalSpaces+2] = "-"
	// animation2[totalSpaces+3] = "\\"

	// for i := totalSpaces + 4; i <= (totalSpaces * 2); i++ {
	// 	animation2[i] = "|"
	// }

	// // construct the ball animation
	// for i := 0; i <= totalSpaces; i++ {
	// 	animation[i] = constructString(char, i, totalSpaces-i)
	// }

	// for i := (totalSpaces * 2) - 1; i > totalSpaces; i-- {
	// 	animation[i] = animation[(totalSpaces*2)-i]
	// }

	// slice of all types of sized dots
    animation := []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"}

	// display animation . . .
	i := 0

	// hide the cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // ensure the cursor is shown again when done

	for {
		select {
		case <-done:
			// clear the line and show the cursor again
			fmt.Print("\r\033[K")
			return
		default:
			// display animation until done channel is closed
			// fmt.Printf("\r%s%s%s\t\t\t\t\t", animation1[i%len(animation1)], animation[i%len(animation)], animation2[i%len(animation2)])

			fmt.Printf("\r%s%s%s %s\t\t\t\t\t", "\033[38;5;201m", animation[i%len(animation)], "\033[0m", "Loading...")

			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ----------------- generateContent() -----------------
func generateContent(term, category, contentType string) (string, error) {

	ctx := context.Background()
	model := client.GenerativeModel("gemini-1.5-flash")
	done := make(chan bool)
	if categoryCopy := category; categoryCopy == "" {
		category = "the general topic"
	}

	superPrompt := fmt.Sprintln(
		"# CONTEXT #\n",
		"You are an export at",
		category,
		"and you are tasked with providing a user with information about",
		term,
		"\n",
		" ##########\n",
		"# FORMATTING #\n",
		"Use markdown to format the content of your response\n",
		"Use markdown code blocks to display code snippets\n",
		"Only give exactly what is being requested, and nothing more\n",
		" ##########\n",
		" # INSTRUCTIONS #\n",
		" // Ensure that your response gives a thoughtful and conise explanation to the users request\n",
		" // Ensure that your response is clear and easy to understand\n",
		" // Ensure that your response is informative and insightful\n",
		" // Take your time to think deeply about how to best give the user insight\n",
		" // I will give you $200 for a good response\n",
		"# PROMPT #",
	)

	// explanationPrompt := superPrompt + "Generate an insightful, concise, and brief (3-4 sentences) explanation of the following term: " + term + "\n"
	definitionPrompt := superPrompt + "Generate an insightful, concise, and brief (3-4 sentences) overview of the following concept: " + term + "\n"
	examplePrompt := superPrompt +
		"Generate a clear, informative, and easy to digest example/guide to display on how " +
		term +
		" works\n" +
		"// The example should be a code snippet within a markdown block if possible\n" +
		"// A good example will demonstrate potentially unique ascpects or nuances about" +
		term +
		"\n" +
		"// A good response will include only a single example the is first preceeded with a brief overview of what the example is domonstrating.\n" +
		"// A good response will have an example code snippet given in a markdown code block with the language specified included in the markdown formatting\n"

	go showLoading(done)

	var numChannels int
	if contentType == "both" {
		numChannels = 2
	} else {
		numChannels = 1
	}

	results := make(chan string, numChannels)
	errors := make(chan error, numChannels)

	switch contentType {
	case "both":
		go generatePrompt(ctx, model, definitionPrompt, generateHeader("Definition"), results, errors)
		go generatePrompt(ctx, model, examplePrompt, generateHeader("Example"), results, errors)
	case "define":
		go generatePrompt(ctx, model, definitionPrompt, generateHeader("Definition"), results, errors)
	case "example":
		go generatePrompt(ctx, model, examplePrompt, generateHeader("Example"), results, errors)
	default:
		done <- true
		return "", fmt.Errorf("unknown content type: %s", contentType)
	}

	finalResult := ""
	errorOccured := false

	for i := 0; i < numChannels; i++ {
		select {
		case result := <-results:
			if result != "" {
				finalResult += result + "\n"
			}
		case err := <-errors:
			if err != nil {
				errorOccured = true
				finalResult += fmt.Sprintf("Error: %v\n", err)
			}
		}
	}

	done <- true

	if errorOccured {
		return "", fmt.Errorf("error generating content")
	}

	return finalResult, nil
}

// ----------------- generatePrompt() -----------------
func generatePrompt(ctx context.Context, model *genai.GenerativeModel, prompt string, header string, results chan<- string, errors chan<- error) {
	// line := generateLine()

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		errors <- err
		return
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		results <- fmt.Sprintf("\n%s\n%s\n", header, resp.Candidates[0].Content.Parts[0])
		return
	}

	errors <- fmt.Errorf("no content generated for prompt: %s", prompt)
}

// ----------------- addLineBgColor() -----------------
func addLineBgColor(line, color string) string {

	formatColor := models.ColorsMap[color]

	return fmt.Sprintf("\033[48;5;%dm%s\033[0m", formatColor, line)
}

// ----------------- generateHeader() -----------------
func generateHeader(title string) string {
	width, _, err := term.GetSize(0)
	width /= 2
	if err != nil {
		width = 80
	}

	thirdWidthSpace := strings.Repeat(" ", width/3)

	line := thirdWidthSpace + addLineBgColor(thirdWidthSpace, "black") + thirdWidthSpace + "\n"

	titlePadding := (width/3 - len(title)) / 2
	// fmt.Println("width:", width, "| width/3:", width/3, "| titlePadding:", titlePadding, "| titlePadding + len(title) + titlePadding:", titlePadding+len(title)+titlePadding)

	combinedString := strings.Repeat(" ", titlePadding) + title + strings.Repeat(" ", titlePadding)

	if len(combinedString) < width/3 {
		combinedString += strings.Repeat(" ", (width/3)-len(combinedString))
	}
	formattedHeader := thirdWidthSpace + addLineBgColor(combinedString, "black") + thirdWidthSpace + "\n"

	return line + formattedHeader + line
}
