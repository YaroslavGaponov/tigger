package memstore

import (
	"context"
	"errors"

	"github.com/YaroslavGaponov/tigger/internal/store"
)

var (
	ErrDataNotFound = errors.New("data is not found")
)

type MemStore struct {
	data map[string][]byte
}

func NewMemStore() store.Store {
	data := make(map[string][]byte)
	return &MemStore{data: data}
}

func (store *MemStore) Add(ctx context.Context, name string, data []byte) error {
	store.data[name] = data
	return nil
}

func (store *MemStore) Remove(ctx context.Context, name string) error {
	delete(store.data, name)
	return nil
}

func (store *MemStore) Get(ctx context.Context, name string, start, length int) ([]byte, error) {
	data, found := store.data[name]
	if !found {
		return nil, ErrDataNotFound
	}
	return data[start : start+length], nil
}
