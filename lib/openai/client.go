package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	OpenAIBaseURL           = "https://api.openai.com"
	ChatCompletionsEndpoint = "/v1/chat/completions"
)

func sendApiRequest(ctx context.Context, chatRequest ChatRequest) (ChatMessage, error) {
	endPoint := OpenAIBaseURL + ChatCompletionsEndpoint

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return ChatMessage{}, errors.New("Cannot find OPENAI_API_KEY environment variable.")
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("Error while marhsalling JSON : %w", err)
	}

	// create http request
	req, err := http.NewRequestWithContext(ctx, "POST", endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return ChatMessage{}, fmt.Errorf("Error while creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// send the http request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("Reading the HTTP response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ChatMessage{}, fmt.Errorf("End-point returned code %d: %s", resp.StatusCode, string(body))
	}

	// parse the response
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return ChatMessage{}, fmt.Errorf("Error unmarhslling the response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return ChatMessage{}, errors.New("No response from LLM")
	}

	return chatResp.Choices[0].Message, nil
}
