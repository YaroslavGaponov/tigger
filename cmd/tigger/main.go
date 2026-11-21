package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
	"github.com/YaroslavGaponov/tigger/internal/watcher"
)

func main() {
	ctx := context.Background()

	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "Ollama URL")
	ollamaChatModel := flag.String("ollama-chat-model", "qwen3.8:27b", "Ollama chat model")
	ollamaEmbeddingModel := flag.String("ollama-embedding-model", "nomic-embed-text", "Ollama embedding model")
	indexerChunkSize := flag.Int("indexer-chunk-size", 500, "Indexer chunk size")
	indexerChunkOverlap := flag.Int("indexer-chunk-overlap", 100, "Indexer chunk overlap")
	oracleTopKEmbeddings := flag.Int("oracle-topk-embeddings", 10, "Oracle topK embeddings")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] directories...\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	store := store.NewMemory()

	vectordb := vectordb.NewMemory(
		vectordb.MemoryWithDistanceFunc(vectordb.CosineSimilarity),
	)

	ollama := ollama.NewOllama(
		ollama.WithUrl(*ollamaUrl),
		ollama.WithChatModel(*ollamaChatModel),
		ollama.WithEmbeddingModel(*ollamaEmbeddingModel),
	)

	indexer := indexer.NewIndexer(
		indexer.WithChunkSize(*indexerChunkSize),
		indexer.WithChunkOverlap(*indexerChunkOverlap),
	)

	oracle := oracle.NewOracle(
		ollama,
		store,
		vectordb,
		indexer,
		oracle.WithTopKEmbeddings(*oracleTopKEmbeddings),
	)

	watcher, err := watcher.NewFileWatcher()
	if err != nil {
		panic(err)
	}

	watcher.OnCreate(func(fileName string) {
		data, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
		if err := oracle.AddData(ctx, fileName, string(data)); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	watcher.OnRemove(func(fileName string) {
		if err := oracle.RemoveData(ctx, fileName); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	})

	for _, dirName := range flag.Args() {
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

			if err := oracle.AddData(ctx, path, string(data)); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				return
			}
			fmt.Println("done")

		}
	}

	watcher.Watch(ctx, flag.Args())

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("❓")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
			}
			break
		}
		question := scanner.Text()
		if question == "exit" {
			os.Exit(0)
		}
		answer, err := oracle.Ask(ctx, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			continue
		}

		fmt.Println("🚩", answer)
	}

}
