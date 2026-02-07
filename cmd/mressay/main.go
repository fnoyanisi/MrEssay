package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/fnoyanisi/MrEssay/lib/openai"
)

func main() {

	// infrastructure
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Printf("Need 'OPENAI_API_KEY' environemnt variable to use OpenAI API")
		return
	}

	// agents
	assignerRole := "You are a Year 6 primary school teacher in New Zealand. Create a " +
		"one-page writing assignment using the TIDE framework (Topic sentence, " +
		"Important ideas, Detailed explanations, Ending). Use age-appropriate " +
		"language and clear instructions. Provide only the writing topic and a " +
		"few brief tips for starting. No emojis, introductions, or follow-up."
	assignerAgent := openai.NewAgent(ctx, openai.ModelGPT4oMini, apiKey, assignerRole)
	assignment, err := assignerAgent.Chat("Give a topic.")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	reviewerRole := "You are a primary school teacher in New Zealand reviewing " +
		"Year 6 writing assignments. Assess the work at an age-appropriate level " +
		"using the TIDE (Topic sentence, Important ideas, Detailed explanations" +
		" and Ending) framework. Identify strengths, highlight areas for improvement, " +
		"and clearly point out spelling, grammar, punctuation, and typing errors. " +
		"Provide constructive, encouraging feedback."
	reviewerAgent := openai.NewAgent(ctx, openai.ModelGPT4oMini, apiKey, reviewerRole)

	ocrAgent := openai.NewAgent(ctx, openai.ModelGPT4oMini, apiKey, "")
	ocrPrompt := "You are performing strict optical character recognition (OCR)." +
		"The image contains a letter or an essay." +
		"Transcribe exactly the text visible in the image." +
		"- Preserve original spelling, grammar, typos, punctuation, capitalization, and line breaks." +
		"- Do not correct or normalize anything." +
		"- Do not infer missing text." +
		"- Do not add explanations, summaries, or context." +
		"- Do not prepend or append any text." +
		"Output only the raw transcribed text, exactly as it appears."
	img, err := os.ReadFile("./img.jpeg")
	if err != nil {
		fmt.Printf("Error reading the image: %w", err)
		return
	}
	r, err := ocrAgent.ChatWithImage(ocrPrompt, img)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	prompt := "Here is a text about summer holidays from a student. Review it and let me know the outcome." +
		"Giving the text : " + r

	response, err := reviewerAgent.Chat(prompt)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Println("Review response:" + response)

}
