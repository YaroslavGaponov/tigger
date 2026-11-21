package store

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrDataNotFound = errors.New("data is not found")
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemory() Store {
	data := make(map[string]string)
	return &Memory{data: data}
}

func (store *Memory) Add(ctx context.Context, name string, data string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[name] = data
	return nil
}

func (store *Memory) Remove(ctx context.Context, name string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, name)
	return nil
}

func (store *Memory) Get(ctx context.Context, name string, start, length int) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	data, found := store.data[name]
	if !found {
		return "", ErrDataNotFound
	}
	return data[start : start+length], nil
}
