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

	fileName := flag.String("filename", "assets/test.txt", "Filename")
	flag.Parse()

	ollama := ollama.NewOllama()
	store := store.NewInMemoryStore()
	vectordb := vectordb.NewInMemoryVectordb()
	indexer := indexer.NewIndexer()
	oracle := oracle.NewOracle(ollama, store, vectordb, indexer)

	fmt.Print("🚩 loading and indexing...")
	data, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}
	if err := oracle.AddData(ctx, *fileName, string(data)); err != nil {
		panic(err)
	}
	fmt.Println("done.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("❓ question: ")
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
		}

		fmt.Println("🙈 answer: " + answer)

	}

}
