package llm

import (
	"errors"
	"os"
)

type LLMProvider interface {
	GenerateResponse(prompt string, previousId string) (string, string, error)
	GenerateEmbedding(text string) ([]float32, error)
}

func NewLLMProvider() (LLMProvider, error) {
	provider := os.Getenv("LLM_PROVIDER")

	switch provider {
	case "openai":
		return NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_API_MODEL")), nil
	case "lmstudio":
		return NewLmStudioProvider(os.Getenv("LM_STUDIO_MODEL"), os.Getenv("LM_STUDIO_URL")), nil
	default:
		return nil, errors.New("invalid LLM provider")
	}
}

type MockLLM struct{}

func (m *MockLLM) GenerateResponse(prompt string, previousId string) (string, string, error) {
	return "This is a mock response from the LLM.", "", nil
}

func (m *MockLLM) GenerateEmbedding(text string) ([]float32, error) {
	return make([]float32, 1536), nil // Return empty vector of size 1536
}
