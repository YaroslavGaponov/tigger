package store

import "context"

type Store interface {
	Add(ctx context.Context, name string, data []byte) error
	Remove(ctx context.Context, name string) error
	Get(ctx context.Context, name string, start, length int) ([]byte, error)
}
