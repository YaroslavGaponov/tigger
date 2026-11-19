package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/ollama"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
)

func main() {
	ctx := context.Background()

	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "Ollama URL")
	ollamaChatModel := flag.String("ollama-chat-model", "qwen3.8:27b", "Ollama chat model")
	ollamaEmbeddingModel := flag.String("ollama-embedding-model", "nomic-embed-text", "Ollama embedding model")
	indexerChunkSize := flag.Int("indexer-chunk-size", 500, "Indexer chunk size")
	indexerChunkOverlap := flag.Int("indexer-chunk-overlap", 100, "Indexer chunk overlap")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <files...>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	store := store.NewInMemoryStore()

	vectordb := vectordb.NewInMemoryVectordb()

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
	)

	for _, fileName := range flag.Args() {
		fmt.Printf("📂 Loading and indexing file: %s...", fileName)
		data, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		if err := oracle.AddData(ctx, fileName, string(data)); err != nil {
			panic(err)
		}
		fmt.Println("done")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("❓")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input:", err)
			}
			break
		}
		question := scanner.Text()
		if question == "exit" {
			os.Exit(0)
		}
		answer, err := oracle.Ask(ctx, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error oracle ask:", err)
			continue
		}

		fmt.Println("🚩", answer)
	}

}
