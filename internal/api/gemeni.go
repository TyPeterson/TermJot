package api

import (
	"context"
	"fmt"
    "os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var client *genai.Client

// ----------------- InitializeGeminiClient() -----------------
func InitializeGeminiClient() {

    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        fmt.Println("Error: GEMINI_API_KEY environment variable not set")
        return
    }

	ctx := context.Background()
	var clientErr error
	client, clientErr = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if clientErr != nil {
        fmt.Println("Error initializing Gemini client:", clientErr)
	}
}


// ----------------- GenerateDefinition() -----------------
func GenerateDefinition(termName, categoryName string) string {

    var promptCategory string

    if categoryName == "" {
        promptCategory = termName
    } else {
        promptCategory = categoryName
    }

    prompt :=
        "# CONTEXT #\n" +
        "You are an expert at " + promptCategory + " and you are tasked with providing a user with information about " + termName + "\n" +
        "The user is requesting a brief, concise, and insightful definition and overview of " + termName + "\n" +
        " ##########\n" +
        "# INSTRUCTIONS #\n" +
        " // Ensure that your response gives a thoughtful and concise definition/overview/summary of " + termName + "\n" +
        " // Ensure that your response is clear and easy to understand\n" +
        " // Ensure that your response is informative and insightful\n" +
        " // Take your time to think deeply about how to best give the user insight\n" +
        " // Assume not prior knowledge of the topic and so be sure to explain the topic in a way that is easy to understand\n" +
        " // I will give you $200 for a good response\n" +
        "# PROMPT #\n" +
        "Generate an insightful, concise, and brief (3-4 sentences) overview of the following concept: " + termName + "\n"


    return generateContent(prompt)
}

// ----------------- GenerateExample() -----------------
func GenerateExample(termName, categoryName string) string {

    var promptCategory string

    if categoryName == "" {
        promptCategory = termName
    } else {
        promptCategory = categoryName
    }

    prompt :=
        "# CONTEXT #\n" +
        "You are an expert at " + promptCategory + " and you are tasked with providing a user with an example that clearly demonstrates" + termName + "\n" +
        "The user is requesting a clear, informative, and easy to digest example/guide to display how " + termName + " works\n" +
        " ##########\n" +
        "# INSTRUCTIONS #\n" +
        " // Generate a single example that is nested within a markdown code block\n" +
        " // The example should be a code snippet within a markdown code block\n" +
        " // A good example will demonstrate potentially unique aspects or nuances about " + termName + "\n" +
        " // A good response will include only a single example, given in a markdown code block\n" +
        " // A good response will have an example code snippet given in a markdown code block with the language specified included in the markdown formatting\n" +
        " // The example will contain exclusively the code snippet and nothing more\n" +
        " // Ensure that your response gives a thoughtful and concise example/guide to display how " + termName + " works\n" +
        " // Ensure that your response is clear and easy to understand\n" +
        " // Ensure that your response is informative and insightful\n" +
        " // Take your time to think deeply about how to best give the user insight\n" +
        " // I will give you $200 for a good response\n" +
        " ##########\n" +
        "# FORMATTING #\n" +
        "Use markdown code blocks to display code snippets\n" +
        "Only give exactly what is being requested, and nothing more\n" +
        " ##########\n" +
        "# PROMPT #\n" +
        "Generate a clear, informative, and easy to digest example/guide to display on how " + termName + " works\n"


    return generateContent(prompt)
}


// ----------------- generateContent() -----------------
func generateContent(prompt string) string {

	ctx := context.Background()
	model := client.GenerativeModel("gemini-1.5-flash")

    resp, err := model.GenerateContent(ctx, genai.Text(prompt))

    if err != nil {
        fmt.Println("Error generating content:", err)
        return ""
    }

    if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
        return fmt.Sprintf("\n%s\n", resp.Candidates[0].Content.Parts[0])
    }

    return ""
}

