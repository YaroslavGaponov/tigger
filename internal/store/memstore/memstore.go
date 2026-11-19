package memstore

import (
	"context"
	"errors"
	"sync"

	"github.com/YaroslavGaponov/tigger/internal/store"
)

var (
	ErrDataNotFound = errors.New("data is not found")
)

type MemStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMemStore() store.Store {
	data := make(map[string][]byte)
	return &MemStore{data: data}
}

func (store *MemStore) Add(ctx context.Context, name string, data []byte) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[name] = data
	return nil
}

func (store *MemStore) Remove(ctx context.Context, name string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, name)
	return nil
}

func (store *MemStore) Get(ctx context.Context, name string, start, length int) ([]byte, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	data, found := store.data[name]
	if !found {
		return nil, ErrDataNotFound
	}
	return data[start : start+length], nil
}
