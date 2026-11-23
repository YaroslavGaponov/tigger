package oracle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YaroslavGaponov/tigger/internal/config"
	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
	"github.com/YaroslavGaponov/tigger/internal/watcher"
)

type Oracle struct {
	ollama   *ollama.Ollama
	store    store.Store
	vectordb vectordb.VectorDB
	indexer  *indexer.Indexer
	watcher  *watcher.FileWatcher

	topKEmbeddings int
	directories    []string
}

func NewOracle(config config.Config) (*Oracle, error) {
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

	watcher, err := watcher.NewFileWatcher()
	if err != nil {
		return nil, err
	}

	oracle := &Oracle{
		store:          store,
		vectordb:       vectordb,
		ollama:         ollama,
		indexer:        indexer,
		watcher:        watcher,
		topKEmbeddings: config.OracleTopKEmbeddings,
		directories:    config.Directories,
	}

	return oracle, nil
}

func (o *Oracle) AddData(ctx context.Context, name, text string) error {
	if err := o.store.Add(ctx, name, text); err != nil {
		return err
	}

	for sentence := range o.indexer.Split(string(text)) {
		vector, err := o.ollama.GetEmbedding(ctx, sentence.Text)
		if err != nil {
			return err
		}
		metadata := map[string]any{
			"name":   name,
			"start":  sentence.Start,
			"length": sentence.Length,
		}
		o.vectordb.Add(ctx, vector, metadata)
	}
	return nil
}

func (o *Oracle) RemoveData(ctx context.Context, name string) error {
	metadata := map[string]any{"name": name}
	if err := o.vectordb.Remove(ctx, metadata); err != nil {
		return err
	}
	return o.store.Remove(ctx, name)

}

func (o *Oracle) Ask(ctx context.Context, question string) (string, error) {
	vector, err := o.ollama.GetEmbedding(ctx, question)
	if err != nil {
		return "", err
	}

	results, err := o.vectordb.Search(ctx, vector, o.topKEmbeddings)
	if err != nil {
		return "", err
	}
	var context []string
	for _, result := range results {
		data, err := o.store.Get(ctx, result["name"].(string), result["start"].(int), result["length"].(int))
		if err != nil {
			return "", err
		}
		context = append(context, string(data))
	}

	messages := []ollama.Message{
		{
			Role:    "system",
			Content: "You are an assistant that answers questions exclusively based on the provided context. If the answer cannot be found in the context, respond: \"Information not found in the text.\" Do not use external knowledge.",
		},
		{
			Role:    "system",
			Content: "Answer must be only one sentence.",
		},
		{
			Role:    "user",
			Content: "Context: " + strings.Join(context, "\n") + "\n",
		}, {
			Role:    "user",
			Content: "Question: " + question,
		},
	}

	answer, err := o.ollama.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (s *Oracle) Run(ctx context.Context) error {

	for _, dirName := range s.directories {
		entries, err := os.ReadDir(dirName)
		if err != nil {
			return err
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

			if err := s.AddData(ctx, path, string(data)); err != nil {
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
		if err := s.AddData(ctx, fileName, string(data)); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	s.watcher.OnRemove(func(fileName string) {
		if err := s.RemoveData(ctx, fileName); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	return s.watcher.Watch(ctx, s.directories)
}
