package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/YaroslavGaponov/tigger/internal/indexer"
	"github.com/YaroslavGaponov/tigger/internal/llm/ollama"
	"github.com/YaroslavGaponov/tigger/internal/store"
	"github.com/YaroslavGaponov/tigger/internal/vectordb"
)

func main() {
	ctx := context.Background()

	fileName := flag.String("filename", "assets/test.txt", "Filename")
	flag.Parse()

	fmt.Printf("🚩loading %s file...", *fileName)
	text, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("done")

	store := store.NewInMemoryStore()
	store.Add(ctx, *fileName, text)

	fmt.Printf("🚩indexing %s file...", *fileName)
	llm := ollama.NewOllama()
	indexer := indexer.NewIndexer()
	vectorDb := vectordb.NewInMemoryVectordb()
	for sentence := range indexer.Split(string(text)) {
		vector, err := llm.GetEmbedding(ctx, sentence.Text)
		if err != nil {
			panic(err)
		}
		metadata := map[string]any{
			"name":   *fileName,
			"start":  sentence.Start,
			"length": sentence.Length,
		}
		vectorDb.Add(ctx, vector, metadata)
	}
	fmt.Println("done")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("❓question: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input:", err)
			}
			break
		}
		question := scanner.Text()

		vector, err := llm.GetEmbedding(ctx, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error llm get embedding:", err)
			continue
		}

		results, err := vectorDb.Search(ctx, vector, 5)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error verctordb search:", err)
			continue
		}
		var context []string
		for _, result := range results {
			data, err := store.Get(ctx, result["name"].(string), result["start"].(int), result["length"].(int))
			if err != nil {
				fmt.Fprintln(os.Stderr, "error store get:", err)
				continue
			}
			context = append(context, string(data))
		}

		answer, err := llm.Chat(ctx, context, question)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error llm chat:", err)
			continue
		}
		fmt.Println("🙈 answer: " + answer)

	}

}
