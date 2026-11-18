package llm

import "context"

type LLM interface {
	GetEmbedding(ctx context.Context, text string) ([]float64, error)
	Chat(ctx context.Context, context []string, question string) (string, error)
}
