package llm

import (
	"context"
	"fmt"
	"time"
)

type LLMResponse struct {
	Response  string             `json:"response"`
	Timestamp time.Time          `json:"timestamp"`
}

type LLMClient interface {
	GenerateResponse(ctx context.Context, text string) (*LLMResponse, error)
	Close() error
}	

type LLMProvider string

const (
	LLMProviderOllama LLMProvider = "ollama"
)

func NewLLMClient(provider LLMProvider) (LLMClient, error) {
	switch provider {
	case LLMProviderOllama:
		ollamaHost := "http://localhost:11434"
		model := "llama3.2:3b"
		return NewOllamaLLMClient(ollamaHost, model)
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
