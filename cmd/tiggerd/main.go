package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/YaroslavGaponov/tigger/internal/api"
	"github.com/YaroslavGaponov/tigger/internal/config"
	"github.com/YaroslavGaponov/tigger/internal/logger"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
	"github.com/YaroslavGaponov/tigger/internal/server"
	"github.com/YaroslavGaponov/tigger/internal/utils"
)

func main() {

	ctx := context.Background()

	config := config.Config{
		Directories:          strings.Split(utils.Getenv("DIRECTORIES", "assets"), ","),
		OllamaUrl:            utils.Getenv("OLLAMA_URL", "http://localhost:11434"),
		OllamaChatModel:      utils.Getenv("OLLAMA_CHAT_MODEL", "qwen3.8:27b"),
		OllamaEmbeddingModel: utils.Getenv("OLLAMA_EMBEDDING_MODEL", "nomic-embed-text"),
		IndexerChunkSize:     utils.Getenv("INDEXER_CHUNK_SIZE", 500),
		IndexerChunkOverlap:  utils.Getenv("INDEXER_CHUNK_OVERLAP", 100),
		OracleTopKEmbeddings: utils.Getenv("ORACLE_TOPK_EMBEDDINGS", 10),
	}

	oracle, err := oracle.NewOracle(config)
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(oracle)
	if err != nil {
		panic(err)
	}

	if err := server.Run(ctx); err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Handler: api.Handler(server),
		Addr:    net.JoinHostPort("0.0.0.0", utils.Getenv("PORT", "7878")),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logger.Log.Info().Msg("👍 HTTP server is starting....")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(2)
		}
	}()

	<-stop

	logger.Log.Info().Msg("👎 Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		os.Exit(2)
	}
}
