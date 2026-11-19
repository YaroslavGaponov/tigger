package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrEmptyEmbeddings = errors.New("empty embeddings")
)

type Ollama struct {
	client         *http.Client
	url            string
	chatModel      string
	embeddingModel string
}

func NewOllama(opts ...OllamaOption) *Ollama {
	ollama := &Ollama{
		client:         &http.Client{Timeout: 10 * time.Second},
		url:            "http://localhost:11434",
		chatModel:      "qwen3.8:27b",
		embeddingModel: "nomic-embed-text",
	}

	for _, opt := range opts {
		opt(ollama)
	}
	return ollama
}

func (ollama *Ollama) GetEmbedding(ctx context.Context, text string) ([]float64, error) {

	reqBody := embedRequest{
		Model: ollama.embeddingModel,
		Input: text,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollama.url+"/api/embed", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := ollama.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, ErrEmptyEmbeddings
	}

	return result.Embeddings[0], nil

}

func (llm *Ollama) Chat(ctx context.Context, messages []Message) (string, error) {

	request := ChatRequest{
		Model:    llm.chatModel,
		Messages: messages,
		Stream:   false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		llm.url+"/api/chat",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(data))
	}

	var result ChatResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
