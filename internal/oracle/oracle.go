package oracle

import (
	"context"
	"strings"

	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
)

type Oracle struct {
	ollama   *ollama.Ollama
	store    store.Store
	vectordb vectordb.VectorDB
	indexer  *indexer.Indexer
}

func NewOracle(ollama *ollama.Ollama, store store.Store, vectordb vectordb.VectorDB, indexer *indexer.Indexer) *Oracle {
	return &Oracle{
		ollama:   ollama,
		store:    store,
		vectordb: vectordb,
		indexer:  indexer,
	}
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

func (o *Oracle) RemovData(ctx context.Context, name string) error {
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

	results, err := o.vectordb.Search(ctx, vector, 5)
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
