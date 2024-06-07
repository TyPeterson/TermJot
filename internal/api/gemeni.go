package api

import (
	"context"
	"fmt"
	"log"
	"github.com/google/generative-ai-go/genai"
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

    apiKey := "AIzaSyCrFqdIPZnSsvJJvsL7vcUe7weTFehnGLQ"
	// apiKey := os.Getenv("GEMINI_API_KEY")
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



// ----------------- GenerateDefinition() -----------------
func GenerateDefinition(term, category string) string {
	return generateContent(term, category, "define")
}

// ----------------- GenerateExample() -----------------
func GenerateExample(term, category string) string {
	return generateContent(term, category, "example")
}


// ----------------- generateContent() -----------------
func generateContent(term, category, responseType string) string {

	ctx := context.Background()
	model := client.GenerativeModel("gemini-1.5-flash")
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

	definitionPrompt := superPrompt + "Generate an insightful, concise, and brief (3-4 sentences) overview of the following concept: " + term + "\n" + "DO NOT USE ANY MARKDOWN FORMATTING IN THE RESPONSE\n"
	examplePrompt := superPrompt +
		"Generate a clear, informative, and easy to digest example/guide to display on how " +
		term +
		" works\n" +
		"// The example should be a code snippet within a markdown code block\n" +
		"// A good example will demonstrate potentially unique ascpects or nuances about" +
		term +
		"\n" +
		"// A good response will include only a single example, given in a markdown code block\n" +
		"// A good response will have an example code snippet given in a markdown code block with the language specified included in the markdown formatting\n" +
        "// The example will contain exclusively the code snippet and nothing more\n"


    results := ""

	switch responseType {
        case "define":
            results = generatePrompt(ctx, model, definitionPrompt)
        case "example":
            results = generatePrompt(ctx, model, examplePrompt)
        default:
            return ""
	}


	finalResult := results

	return finalResult
}

// ----------------- generatePrompt() -----------------
func generatePrompt(ctx context.Context, model *genai.GenerativeModel, prompt string) string {

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


