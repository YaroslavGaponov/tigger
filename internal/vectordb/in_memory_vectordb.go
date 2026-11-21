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
	mu             sync.RWMutex
	data           map[*Node]struct{}
	distanceMetric DistanceMetric
}

func NewInMemoryVectordb() VectorDB {
	return &InMemoryVectordb{
		data:           make(map[*Node]struct{}),
		distanceMetric: CosineSimilarity,
	}
}

func (db *InMemoryVectordb) Add(ctx context.Context, embedding []float64, metadata map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	node := &Node{
		embedding: embedding,
		metadata:  metadata,
	}
	db.data[node] = struct{}{}
	return nil
}

func (db *InMemoryVectordb) Remove(ctx context.Context, metadata map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	for node := range db.data {
		mustDelete := true
		for k, v := range metadata {
			if value, found := node.metadata[k]; found && v != value {
				mustDelete = false
				break
			}
		}
		if mustDelete {
			delete(db.data, node)
		}
	}
	return nil
}

func (db *InMemoryVectordb) Search(ctx context.Context, embedding []float64, topK int) ([]map[string]any, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	arr := make([]searchNode, 0, len(db.data))
	for node := range db.data {
		arr = append(arr, searchNode{
			Node:  node,
			Scope: db.distanceMetric(embedding, node.embedding),
		},
		)
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

func (db *InMemoryVectordb) Size(ctx context.Context) int {
	return len(db.data)
}
