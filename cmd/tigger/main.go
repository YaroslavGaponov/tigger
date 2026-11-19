package main

import (
	"bufio"
	"context"
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

	ollama := ollama.NewOllama()
	store := store.NewInMemoryStore()
	vectordb := vectordb.NewInMemoryVectordb()
	indexer := indexer.NewIndexer()
	oracle := oracle.NewOracle(ollama, store, vectordb, indexer)

	for _, fileName := range os.Args[1:] {
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
		answer, err := oracle.Ask(ctx, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error oracle ask:", err)
			continue
		}

		fmt.Println("🚩", answer)

	}

}
