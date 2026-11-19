package vectordb

import (
	"context"
	"sort"
	"sync"
)

type Node struct {
	embedding []float64
	metadata  map[string]any
}

type searchNode struct {
	Node  *Node
	Scope float64
}

type InMemoryVectordb struct {
	mu   sync.RWMutex
	data []Node
}

func NewInMemoryVectordb() VectorDB {
	return &InMemoryVectordb{}
}

func (db *InMemoryVectordb) Add(ctx context.Context, embedding []float64, metadata map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	node := Node{
		embedding: embedding,
		metadata:  metadata,
	}
	db.data = append(db.data, node)
	return nil
}

func (db *InMemoryVectordb) Search(ctx context.Context, embedding []float64, topK int) ([]map[string]any, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	arr := make([]searchNode, len(db.data))
	for i, node := range db.data {
		arr[i] = searchNode{
			Node:  &node,
			Scope: calcDistance(embedding, node.embedding),
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Scope > arr[j].Scope
	})
	var result []map[string]any
	for i := range min(topK, len(arr)) {
		result = append(result, arr[i].Node.metadata)
	}
	return result, nil
}

func calcDistance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}
