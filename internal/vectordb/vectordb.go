package vectordb

import "context"

type DistanceMetric func(a, b []float64) float64

type VectorDB interface {
	Add(ctx context.Context, embedding []float64, metadata map[string]any) error
	Remove(ctx context.Context, metadata map[string]any) error
	Search(ctx context.Context, embedding []float64, topK int) ([]map[string]any, error)
	Size() int
}
