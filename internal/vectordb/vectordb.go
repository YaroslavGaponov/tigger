package vectordb

import "context"

type VectorDB interface {
	Add(ctx context.Context, embedding []float64, metadata map[string]any) error
	Search(ctx context.Context, embedding []float64, topK int) ([]map[string]any, error)
}
