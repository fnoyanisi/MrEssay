package openai

import "context"

type GptModel string

const (
	ModelGPT4oMini GptModel = "gpt-4o-mini"
	ModelGPT5Mini  GptModel = "gpt-5-mini"
)

type Agent struct {
	Model       GptModel
	ApiKey      string
	Role        string
	Context     context.Context
	Temperature float32
	Messages    []ChatMessage
}

func NewAgent(ctx context.Context, model GptModel, apiKey string, role string) *Agent {
	// assign some defaults
	if model == "" {
		model = ModelGPT4oMini
	}

	if role == "" {
		role = "You are a helpful assistant. Answer clearly and concisely. Do not add fluff."
	}

	agent := &Agent{
		Context:     ctx,
		ApiKey:      apiKey,
		Role:        role,
		Model:       model,
		Temperature: 0.0,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: role,
			},
		},
	}

	return agent
}
