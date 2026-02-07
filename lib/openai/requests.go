package openai

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// data structures for API requests / responses

type ImageDetail string

const (
	ImageDetailLow  ImageDetail = "low"
	ImageDetailHigh ImageDetail = "high"
)

type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // either string or []ContentPart
}

type ImageUrl struct {
	URL    string      `json:"url"`
	Detail ImageDetail `json:"detail,omitempty"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageUrl *ImageUrl `json:"image_url,omitempty"`
}

type ChatRequest struct {
	Model       GptModel      `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature *float32      `json:"temperature,omitempty"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
}

// functions

func (a *Agent) AskWithImage(prompt string, img []byte) (string, error) {
	// read the image
	b64 := base64.StdEncoding.EncodeToString(img)
	imgUri := "data:image/jpeg;base64," + b64

	imageUrl := ImageUrl{
		URL:    imgUri,
		Detail: ImageDetailLow,
	}

	// create the payload
	parts := []ContentPart{
		{
			Type: "text",
			Text: prompt,
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

	return a.send(userMessage)
}

func (a *Agent) Ask(message string) (string, error) {
	userMessage := ChatMessage{
		Role:    "user",
		Content: message,
	}

	return a.send(userMessage)
}

// change the temperature before sending the message
func (a *Agent) WithTemperature(temp float32) (*Agent, error) {
	if temp < 0.0 || temp > 2.0 {
		return nil, fmt.Errorf("Invalid temperature value : %f", temp)
	}
	copy := *a
	copy.Temperature = temp
	return &copy, nil
}

// send the message and update the memeroy, not public
func (a *Agent) send(chatMessage ChatMessage) (string, error) {
	a.Messages = append(a.Messages, chatMessage)

	payload := ChatRequest{
		Model:       a.Model,
		Messages:    a.Messages,
		Temperature: &a.Temperature,
	}

	response, err := sendApiRequest(a.Context, payload)
	if err != nil {
		return "", err
	}

	a.Messages = append(a.Messages, response)

	content, ok := response.Content.(string)
	if !ok {
		return "", errors.New("The response content is not a string")
	}

	return content, nil
}
