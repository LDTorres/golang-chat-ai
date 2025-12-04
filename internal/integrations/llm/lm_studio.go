package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gofiber/fiber/v2/log"
)

type LmStudioProvider struct {
	Model   string
	BaseURL string
}

func NewLmStudioProvider(model string, baseURL string) *LmStudioProvider {
	return &LmStudioProvider{
		Model:   model,
		BaseURL: baseURL,
	}
}

func (p *LmStudioProvider) GetModels() ([]string, error) {
	url := p.BaseURL + "/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return nil, err
	}

	models := make([]string, len(result.Data))
	for i, model := range result.Data {
		models[i] = model.ID
	}

	return models, nil
}

func (p *LmStudioProvider) GenerateResponse(prompt string, previousId string) (string, string, error) {
	models, err := p.GetModels()
	if err != nil {
		return "", "", err
	}

	if !slices.Contains(models, p.Model) {
		return "", "", fmt.Errorf("model %s not found", p.Model)
	}

	log.Info("Models: ", models)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": p.Model, // Default, LM Studio might ignore or require specific
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      false,
	})
	if err != nil {
		log.Error(requestBody, err)
		return "", "", err
	}

	url := p.BaseURL + "/chat/completions"

	log.Info("Requesting LLM: ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("API error: ", string(body))
		return "", "", fmt.Errorf("API error: %s", string(body))
	}

	var result struct {
		Id      string `json:"id"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, result.Id, nil
	}

	return "", "", fmt.Errorf("no response from LLM")
}

func (p *LmStudioProvider) GenerateEmbedding(text string) ([]float32, error) {
	// TODO: Implement OpenAI embedding API
	return nil, fmt.Errorf("openai embedding not implemented")
}
