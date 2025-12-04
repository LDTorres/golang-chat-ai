package qdrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type QdrantClient struct {
	BaseURL string
}

func NewQdrantClient(baseURL string) *QdrantClient {
	return &QdrantClient{BaseURL: baseURL}
}

func (c *QdrantClient) CreateCollection(name string, vectorSize int) error {
	url := fmt.Sprintf("%s/collections/%s", c.BaseURL, name)
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create collection: %s", string(b))
	}
	return nil
}

func (c *QdrantClient) UpsertPoints(collectionName string, points []map[string]interface{}) error {
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", c.BaseURL, collectionName)
	body := map[string]interface{}{
		"points": points,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upsert points: %s", string(b))
	}
	return nil
}

func (c *QdrantClient) Search(collectionName string, vector []float32, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/collections/%s/points/search", c.BaseURL, collectionName)
	body := map[string]interface{}{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to search: %s", string(b))
	}

	var result struct {
		Result []map[string]interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}
