package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/YaroslavGaponov/tigger/internal/config"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
)

func main() {
	ctx := context.Background()

	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "Ollama URL")
	ollamaChatModel := flag.String("ollama-chat-model", "qwen3.8:27b", "Ollama chat model")
	ollamaEmbeddingModel := flag.String("ollama-embedding-model", "nomic-embed-text", "Ollama embedding model")
	indexerChunkSize := flag.Int("indexer-chunk-size", 500, "Indexer chunk size")
	indexerChunkOverlap := flag.Int("indexer-chunk-overlap", 100, "Indexer chunk overlap")
	oracleTopKEmbeddings := flag.Int("oracle-topk-embeddings", 10, "Oracle topK embeddings")
	oracleKind := flag.String("oracle-kind", "fast", "Oracle kind")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] directories...\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	config := config.Config{
		Directories:          flag.Args(),
		OllamaUrl:            *ollamaUrl,
		OllamaChatModel:      *ollamaChatModel,
		OllamaEmbeddingModel: *ollamaEmbeddingModel,
		IndexerChunkSize:     *indexerChunkSize,
		IndexerChunkOverlap:  *indexerChunkOverlap,
		OracleTopKEmbeddings: *oracleTopKEmbeddings,
	}

	oracle, err := oracle.NewOracle(config)
	if err != nil {
		panic(err)
	}
	oracle.Run(ctx)

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
		answer, err := oracle.Ask(ctx, *oracleKind, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			continue
		}

		fmt.Println("🚩", answer)
	}

}
