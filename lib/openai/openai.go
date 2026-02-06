// some definitions related to OpenAI API
package openai

import "context"

// consttants

const (
	OpenAIBaseURL           = "https://api.openai.com"
	ChatCompletionsEndpoint = "/v1/chat/completions"
)

type GptModel string

const (
	GPT4oMini GptModel = "gpt-4o-mini"
	GPT5Mini  GptModel = "gpt-5-mini"
)

type ImageDetail string

const (
	ImageDetailLow  ImageDetail = "low"
	ImageDetailHigh ImageDetail = "high"
)

// data structures for API requests / responses

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

type Agent struct {
	Model    GptModel
	ApiKey   string
	Role     string
	Context  context.Context
	Messages []ChatMessage
}

func NewAgent(ctx context.Context, model GptModel, apiKey string, role string) *Agent {
	agent := &Agent{
		Context: ctx,
		ApiKey:  apiKey,
		Role:    role,
	}

	// assign defaults
	if model == "" {
		agent.Model = model
	} else {
		agent.Model = GPT4oMini
	}

	if role == "" {
		agent.Role = "You are a helpful assistant. Answer clearly and concisely. Do not add fluff."
	}

	return agent
}
