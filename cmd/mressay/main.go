package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/fnoyanisi/MrEssay/lib/openai"
)

func main() {

	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Printf("Need 'OPENAI_API_KEY' environemnt variable to use OpenAI API")
		return
	}

	imgToTxtAgent := openai.NewAgent(ctx, openai.GPT4oMini, apiKey, "")

	img, err := os.ReadFile("./img.jpeg")
	if err != nil {
		fmt.Printf("Error reading the image: %w", err)
		return
	}

	r, err := imgToTxtAgent.ImageToText(img)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Printf("Response : %s", r)

}
