// some definitions related to OpenAI API
package main

const (
	OpenAIBaseURL           = "https://api.openai.com"
	ChatCompletionsEndpoint = "/v1/chat/completions"
)

type GptModel string

const (
	GPT4oMini GptModel = "gpt-4o-mini"
	GPT5Mini  GptModel = "gpt-5-mini"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       GptModel      `json:"model"`
	Message     []ChatMessage `json:"messages"`
	Temperature float32       `json:"temperature,omitempty"`
}

type ChatChoice struct {
	Index        int64       `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
}
