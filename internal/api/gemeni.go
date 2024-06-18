package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var client *genai.Client

// ----------------- InitializeGeminiClient() -----------------
func InitializeGeminiClient() {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        fmt.Println("Error: GEMINI_API_KEY not set")
        return
    }

    ctx := context.Background()
    var clientErr error
    client, clientErr = genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if clientErr != nil {
        log.Fatal(clientErr)
    }
}

// ----------------- GenerateResponse() -----------------
func GetResponse(prompt, responseType string) string {
    promptContext := superPrompts["context"]
    promptInstructions := superPrompts[responseType+"_instructions"]
    promptFormatting := superPrompts[responseType+"_formatting"]
    promptExample := superPrompts[responseType+"_examples"]

    completePromp := fmt.Sprintf("%s\n%s\n%s\n%s\n### PROMPT ###\n%s", promptContext, promptInstructions, promptFormatting, promptExample, prompt)

    return generateContent(completePromp)
 }


// ----------------- generateContent() -----------------
func generateContent(prompt string) string {
    InitializeGeminiClient()
    ctx := context.Background()
    model := client.GenerativeModel("gemini-1.5-flash")

    resp, err := model.GenerateContent(ctx, genai.Text(prompt))

    if err != nil {
        log.Fatal(err)
    }

    // fmt.Println("len(resp.Candidates):", len(resp.Candidates))
    if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
        return fmt.Sprintf("\n%s\n", resp.Candidates[0].Content.Parts[0])
    }

    return ""
}
