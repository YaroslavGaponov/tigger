package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/YaroslavGaponov/tigger/internal/config"
	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
	"github.com/YaroslavGaponov/tigger/internal/watcher"
)

type Core struct {
	store    store.Store
	vectordb vectordb.VectorDB
	ollama   *ollama.Ollama
	indexer  *indexer.Indexer
	oracle   *oracle.Oracle
	watcher  *watcher.FileWatcher

	directories []string
}

func NewCore(config config.Config) *Core {
	store := store.NewMemory()

	vectordb := vectordb.NewMemory(
		vectordb.MemoryWithDistanceFunc(vectordb.CosineSimilarity),
	)

	ollama := ollama.NewOllama(
		ollama.WithUrl(config.OllamaUrl),
		ollama.WithChatModel(config.OllamaChatModel),
		ollama.WithEmbeddingModel(config.OllamaEmbeddingModel),
	)

	indexer := indexer.NewIndexer(
		indexer.WithChunkSize(config.IndexerChunkSize),
		indexer.WithChunkOverlap(config.IndexerChunkOverlap),
	)

	oracle := oracle.NewOracle(
		ollama,
		store,
		vectordb,
		indexer,
		oracle.WithTopKEmbeddings(config.OracleTopKEmbeddings),
	)

	watcher, err := watcher.NewFileWatcher()
	if err != nil {
		panic(err)
	}

	return &Core{
		store:       store,
		vectordb:    vectordb,
		ollama:      ollama,
		indexer:     indexer,
		oracle:      oracle,
		watcher:     watcher,
		directories: config.Directories,
	}
}

func (s *Core) Run(ctx context.Context) error {

	for _, dirName := range s.directories {
		entries, err := os.ReadDir(dirName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fmt.Printf("📂 Loading and indexing file: %s...", entry.Name())

			path := filepath.Join(dirName, entry.Name())

			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}

			if err := s.oracle.AddData(ctx, path, string(data)); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				return err
			}
			fmt.Println("done")
		}

	}

	s.watcher.OnCreate(func(fileName string) {
		data, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
		if err := s.oracle.AddData(ctx, fileName, string(data)); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	s.watcher.OnRemove(func(fileName string) {
		if err := s.oracle.RemoveData(ctx, fileName); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	return s.watcher.Watch(ctx, s.directories)
}

func (s *Core) Ask(ctx context.Context, question string) (string, error) {
	return s.oracle.Ask(ctx, question)
}
