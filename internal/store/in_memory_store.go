package store

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrDataNotFound = errors.New("data is not found")
)

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewInMemoryStore() Store {
	data := make(map[string][]byte)
	return &InMemoryStore{data: data}
}

func (store *InMemoryStore) Add(ctx context.Context, name string, data []byte) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[name] = data
	return nil
}

func (store *InMemoryStore) Remove(ctx context.Context, name string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, name)
	return nil
}

func (store *InMemoryStore) Get(ctx context.Context, name string, start, length int) ([]byte, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	data, found := store.data[name]
	if !found {
		return nil, ErrDataNotFound
	}
	return data[start : start+length], nil
}
