package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func sendRequest(ctx context.Context, chatRequest ChatRequest) (string, error) {
	endPoint := OpenAIBaseURL + ChatCompletionsEndpoint

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("Cannot find OPENAI_API_KEY environment variable.")
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return "", fmt.Errorf("Error while marhsalling JSON : %w", err)
	}

	// create http request
	req, err := http.NewRequestWithContext(ctx, "POST", endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error while creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// send the http request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Reading the HTTP response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("End-point returned code %d: %s", resp.StatusCode, string(body))
	}

	// parse the response
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("Error unmarhslling the response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("No response from LLM")
	}

	content, ok := chatResp.Choices[0].Message.Content.(string)
	if !ok {
		return "", errors.New("The response content is not a string")
	}

	return content, nil
}

func (a *Agent) TextFromImage(img []byte) (string, error) {
	// read the image
	b64 := base64.StdEncoding.EncodeToString(img)
	imgUri := "data:image/jpeg:base64," + b64

	imageUrl := ImageUrl{
		URL:    imgUri,
		Detail: ImageDetailLow,
	}

	// create the payload
	parts := []ContentPart{
		{
			// prompt, used an LLM to generate it ;)
			Type: "text",
			Text: "You are performing strict optical character recognition (OCR)." +
				"The image contains a letter or an essay." +
				"Transcribe exactly the text visible in the image." +
				"- Preserve original spelling, grammar, typos, punctuation, capitalization, and line breaks." +
				"- Do not correct or normalize anything." +
				"- Do not infer missing text." +
				"- Do not add explanations, summaries, or context." +
				"- Do not prepend or append any text." +
				"Output only the raw transcribed text, exactly as it appears.",
		},
		{
			Type:     "image_url",
			ImageUrl: &imageUrl,
		},
	}

	userMessage := ChatMessage{
		Role:    "user",
		Content: parts,
	}

	var messages []ChatMessage
	messages = append(messages, userMessage)

	temp := float32(0.0)
	payload := ChatRequest{
		Model:       GPT4oMini,
		Messages:    messages,
		Temperature: &temp,
	}

	return sendRequest(a.Context, payload)
}

func (a *Agent) TextFromText(message string) (string, error) {

	// create the payload
	parts := []ContentPart{
		{
			Type: "test",
			Text: a.Role,
		},
		{
			Type: "text",
			Text: message,
		},
	}

	userMessage := ChatMessage{
		Role:    "user",
		Content: parts,
	}

	var messages []ChatMessage
	messages = append(messages, userMessage)

	temp := float32(0.0)
	payload := ChatRequest{
		Model:       GPT4oMini,
		Messages:    messages,
		Temperature: &temp,
	}

	return sendRequest(a.Context, payload)
}
