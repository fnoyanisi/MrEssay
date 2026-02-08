package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	openai "github.com/fnoyanisi/MrEssay/lib/openai"
	telegram "github.com/fnoyanisi/MrEssay/lib/telegram"
)

func main() {

	////////////////////////////////////////////////////////
	// infrastructure
	////////////////////////////////////////////////////////
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Printf("Need 'OPENAI_API_KEY' environemnt variable to use OpenAI API")
		return
	}

	tgIDStr := os.Getenv("TELEGRAM_ID")
	if tgIDStr == "" {
		fmt.Println("Cannot find TELEGRAM_ID environment variable")
		return
	}
	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Error while parsing %s : %v", tgIDStr, err)
		return
	}

	tgApiKey := os.Getenv("TELEGRAM_API_TOKEN")
	if apiKey == "" {
		fmt.Println("Cannot find TELEGRAM_API_TOKEN environment variable")
		return
	}

	tgb, err := telegram.NewTelegramBot(tgApiKey, tgID)
	if err != nil {
		fmt.Errorf("Error : %v", err)
		return
	}

	////////////////////////////////////////////////////////
	// end of infrastructure
	////////////////////////////////////////////////////////

	greeting := fmt.Sprintf("Hi there! Today is %s.", time.Now().Format("02 January, 2006"))
	if err := tgb.SendMessage(greeting); err != nil {
		fmt.Errorf("Error : %v", err)
		return
	}

	////////////////////////////////////////////////////////
	// send the assignment
	////////////////////////////////////////////////////////
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
	if err := tgb.SendMessage(assignment); err != nil {
		fmt.Errorf("Error : %v", err)
		return
	}

	if err := tgb.SendMessage("Send me a photo of your writing when you finish it."); err != nil {
		fmt.Errorf("Error : %v", err)
		return
	}

	// wait for the user
	// this is blocking
	var path string
	for update := range tgb.Listen() {
		if update.Message.Chat.ID != tgb.GetUserId() {
			continue
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Photo != nil {
			// not the best way, but...
			os.MkdirAll("downloads", 0755)
			currentTime := time.Now()
			filename := fmt.Sprintf("photo_%s.jpg", currentTime.Format("20060102150405"))
			path = filepath.Join("downloads", filename)
			if err := tgb.SaveTelegramImage(update, path); err != nil {
				fmt.Errorf("Error : %v", err)
				return
			}

			tgb.SendMessage("Photos has been received. Thank you!")
			tgb.StopListening()
			break
		}
	}

	////////////////////////////////////////////////////////
	// OCR
	////////////////////////////////////////////////////////
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
	img, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading the image: %w", err)
		return
	}

	essayText, err := ocrAgent.ChatWithImage(ocrPrompt, img)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	////////////////////////////////////////////////////////
	// review the essay
	////////////////////////////////////////////////////////
	reviewerRole := "You are a primary school teacher in New Zealand reviewing " +
		"Year 6 writing assignments. Assess the work at an age-appropriate level " +
		"using the TIDE (Topic sentence, Important ideas, Detailed explanations" +
		" and Ending) framework. Identify strengths, highlight areas for improvement, " +
		"and clearly point out spelling, grammar, punctuation, and typing errors. " +
		"Provide constructive, encouraging and concise feedback. Keep the tone appropriate" +
		"for a year 6 student."
	reviewerAgent := openai.NewAgent(ctx, openai.ModelGPT4oMini, apiKey, reviewerRole)

	prompt := "Here is a text about summer holidays from a student. Review it and let me know the outcome." +
		"Giving the text : " + essayText

	reviewResponse, err := reviewerAgent.Chat(prompt)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	tgb.SendMessage(reviewResponse)
}
