package config

type Config struct {
	Directories []string

	OllamaUrl            string
	OllamaChatModel      string
	OllamaEmbeddingModel string

	IndexerChunkSize    int
	IndexerChunkOverlap int

	OracleTopKEmbeddings int
}
