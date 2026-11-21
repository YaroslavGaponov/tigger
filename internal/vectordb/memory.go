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

type Memory struct {
	mu             sync.RWMutex
	data           []Node
	distanceMetric DistanceMetric
}

func NewMemory(opts ...MemoryOption) VectorDB {
	memory := &Memory{
		distanceMetric: CosineSimilarity,
	}

	for _, opt := range opts {
		opt(memory)
	}

	return memory
}

func (db *Memory) Add(ctx context.Context, embedding []float64, metadata map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	node := Node{
		embedding: embedding,
		metadata:  metadata,
	}
	db.data = append(db.data, node)
	return nil
}

func (db *Memory) Remove(ctx context.Context, metadata map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	i, end := 0, len(db.data)
	for i < end {
		node := db.data[i]
		mustDelete := true
		for k, v := range metadata {
			if value, found := node.metadata[k]; found && v != value {
				mustDelete = false
				break
			}
		}
		if mustDelete {
			end--
			db.data[i] = db.data[end]
			db.data[end] = Node{}
		} else {
			i++
		}
	}
	db.data = db.data[:end]
	return nil
}

func (db *Memory) Search(ctx context.Context, embedding []float64, topK int) ([]map[string]any, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	arr := make([]searchNode, 0, len(db.data))
	for _, node := range db.data {
		arr = append(arr, searchNode{
			Node:  &node,
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

func (db *Memory) Size() int {
	return len(db.data)
}
