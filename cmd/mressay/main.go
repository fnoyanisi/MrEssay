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

	}

	img, err := os.ReadFile("./img.jpeg")
	if err != nil {
		fmt.Printf("Error reading the image: %w", err)
		return
	}

	r, err := openai.TextFromImage(ctx, apiKey, img)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Printf("Resonse : %s", r)

}
