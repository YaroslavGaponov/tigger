package oracle

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/YaroslavGaponov/tigger/internal/config"
	filereader "github.com/YaroslavGaponov/tigger/internal/file-reader"
	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/logger"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
	"github.com/YaroslavGaponov/tigger/internal/watcher"
)

const (
	AskKindFast     = "fast"
	AskModeAccurate = "accurate"
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

func (o *Oracle) Ask(ctx context.Context, kind string, question string) (string, error) {
	switch kind {
	case AskModeAccurate:
		return o.accurateAsk(ctx, question)
	case AskKindFast:
		return o.fastAsk(ctx, question)
	default:
		return o.fastAsk(ctx, question)
	}
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

			logger.Log.Info().Msgf("📂 Loading and indexing file: %s...", entry.Name())
			fileName := filepath.Join(dirName, entry.Name())
			context, err := filereader.GetContext(fileName)
			if err != nil {
				logger.Log.Err(err).Str("component", "oracle").Send()
				return err
			}
			if err := s.AddData(ctx, fileName, context); err != nil {
				logger.Log.Err(err).Str("component", "oracle").Send()
				return err
			}
		}
		logger.Log.Info().Msgf("👓 The storage contains %d documents, while the VectorDB contains %d vectors.", s.store.Size(ctx), s.vectordb.Size(ctx))
	}

	s.watcher.OnCreate(func(fileName string) {
		context, err := filereader.GetContext(fileName)
		if err != nil {
			logger.Log.Err(err).Str("component", "oracle").Send()
			return
		}
		if err := s.AddData(ctx, fileName, context); err != nil {
			logger.Log.Err(err).Str("component", "oracle").Send()
			return
		}
		logger.Log.Info().Msgf("👓 The storage contains %d documents, while the VectorDB contains %d vectors.", s.store.Size(ctx), s.vectordb.Size(ctx))
	})

	s.watcher.OnRemove(func(fileName string) {
		if err := s.RemoveData(ctx, fileName); err != nil {
			logger.Log.Err(err).Str("component", "oracle").Send()
			return
		}
		logger.Log.Info().Msgf("👓 The storage contains %d documents, while the VectorDB contains %d vectors.", s.store.Size(ctx), s.vectordb.Size(ctx))
	})

	return s.watcher.Watch(ctx, s.directories)
}

func (o *Oracle) fastAsk(ctx context.Context, question string) (string, error) {
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
		AssistantPrompt,
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

func (o *Oracle) accurateAsk(ctx context.Context, question string) (string, error) {
	vector, err := o.ollama.GetEmbedding(ctx, question)
	if err != nil {
		return "", err
	}

	results, err := o.vectordb.Search(ctx, vector, 5*o.topKEmbeddings)
	if err != nil {
		return "", err
	}
	var chunks []string
	for _, result := range results {
		data, err := o.store.Get(ctx, result["name"].(string), result["start"].(int), result["length"].(int))
		if err != nil {
			return "", err
		}
		chunks = append(chunks, string(data))
	}

	scopes := make([][2]int, len(chunks))
	for i, chunk := range chunks {
		messages := []ollama.Message{
			EstimatePrompt,
			{
				Role:    "user",
				Content: "Question: " + question,
			},
			{
				Role:    "user",
				Content: "Text: " + chunk,
			},
		}
		result, err := o.ollama.Chat(ctx, messages)
		if err != nil {
			return "", err
		}
		scope, err := strconv.Atoi(result)
		if err != nil {
			return "", err
		}
		scopes[i] = [2]int{i, scope}
	}

	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i][1] > scopes[j][1]
	})

	var context []string
	for i := 0; i < min(o.topKEmbeddings, len(scopes)); i++ {
		context = append(context, chunks[scopes[i][0]])
	}

	messages := []ollama.Message{
		AssistantPrompt,
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
