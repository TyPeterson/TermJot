package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var client *genai.Client

// func uwu() {
// 	// fmt.Println("test command running...")
// 	// screenshotCmd := exec.Command("screencapture", "screenshot.png")
// 	// err := screenshotCmd.Run()
// 	// if err != nil {
// 	// 	fmt.Println("Error taking screenshot:", err)
// 	// 	return
// 	// }

// 	// fmt.Println("Screenshot saved as screenshot.png")

// 	ctx := context.Background()
// 	apiKey := "AIzaSyCrFqdIPZnSsvJJvsL7vcUe7weTFehnGLQ"
// 	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	prompt := "tell me about yourself"
// 	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(answer)
// }

func InitializeGeminiClient() {
	apiKey := "AIzaSyCrFqdIPZnSsvJJvsL7vcUe7weTFehnGLQ"
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

func generateContent(prompt string) (string, error) {

	ctx := context.Background()
	model := client.GenerativeModel("gemini-1.5-flash")

	screenshotCmd := exec.Command("screencapture", "screenshot.png")
	err := screenshotCmd.Run()
	if err != nil {
		fmt.Println("Error taking screenshot:", err)
		return "", err
	}

	f, err := os.Open("screenshot.png")
	if err != nil {
		return "", fmt.Errorf("error opening screenshot: %v", err)
	}
	defer f.Close()

	file, err := client.UploadFile(ctx, "", f, nil)

	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	completePrompt :=
		" ### CONTEXT ###\n" +
			"You are an expert programmer and software engineer who is meticulous, thoughtful, and follows best practices.\n" +
			"The following is a user question with an included screenshot of what they are looking at on their screen for context\n" +
			"The terminal that they are typing their question is visible on the screen so ignore the text in the terminal that says: 'termjot test " + prompt + "'\n" +
			" ### INSTRUCTIONS ###\n" +
			"Answer the user question by doing/answering whatever they need, making use of the context of what they can see on their screen\n" +
			"Do as the user requests. Answer questions, give code, help debug, explain concepts, etc.\n" +
			"Take your time to think deeply about the question and solution. Read over everything carefully!\n" +
			"I will give you $200 for a perfect response.\n" +
			" ### QUESTION ###\n" +
			prompt + "\n"

	resp, err := model.GenerateContent(ctx, genai.Text(completePrompt), genai.FileData{URI: file.URI})

	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}

	var finalResult string

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		finalResult = fmt.Sprintf("\n%s\n", resp.Candidates[0].Content.Parts[0])
	}

	return finalResult, nil
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the API",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		InitializeGeminiClient()

		// prompt := "debug my the code I'm looking at"
		content, err := generateContent(prompt)
		if err != nil {
			fmt.Printf("Error generating content: %v\n", err)
			return
		}

		fmt.Println(content)
	},
}
