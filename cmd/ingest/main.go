package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/LDTorres/golang-chat-ai/internal/integrations/llm"
	"github.com/LDTorres/golang-chat-ai/internal/integrations/qdrant"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	qdrantURL := "http://localhost:6333"
	if url := os.Getenv("QDRANT_URL"); url != "" {
		qdrantURL = url
	}

	qdrantClient := qdrant.NewQdrantClient(qdrantURL)
	collectionName := "documents"
	vectorSize := 768 // Nomic embed text v1.5 size

	// Create collection
	err := qdrantClient.CreateCollection(collectionName, vectorSize)
	if err != nil {
		log.Printf("Collection might already exist or error: %v", err)
	}

	// Init LLM
	providerName := os.Getenv("LLM_PROVIDER")
	if providerName == "" {
		providerName = "lmstudio"
	}

	// We need to manually init the provider here since we are not in the main app
	var provider llm.LLMProvider
	switch providerName {
	case "openai":
		provider = llm.NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPEN_API_MODEL"))
	case "lmstudio":
		provider = llm.NewLmStudioProvider(os.Getenv("LM_STUDIO_MODEL"), os.Getenv("LM_STUDIO_URL"))
	default:
		provider = &llm.MockLLM{}
	}

	// Read docs
	docsDir := "./docs"
	err = filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".txt") || strings.HasSuffix(path, ".md")) {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
				return nil
			}

			log.Printf("Processing %s...", path)

			// Simple chunking (e.g., by paragraphs)
			chunks := strings.Split(string(content), "\n\n")
			var points []map[string]interface{}

			for _, chunk := range chunks {
				if strings.TrimSpace(chunk) == "" {
					continue
				}

				embedding, err := provider.GenerateEmbedding(chunk)
				if err != nil {
					log.Printf("Failed to generate embedding for chunk: %v", err)
					continue
				}

				points = append(points, map[string]interface{}{
					"id":     uuid.New().String(),
					"vector": embedding,
					"payload": map[string]interface{}{
						"content": chunk,
						"source":  path,
					},
				})
			}

			if len(points) > 0 {
				if err := qdrantClient.UpsertPoints(collectionName, points); err != nil {
					log.Printf("Failed to upsert points: %v", err)
				} else {
					log.Printf("Successfully ingested %d chunks from %s", len(points), path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
