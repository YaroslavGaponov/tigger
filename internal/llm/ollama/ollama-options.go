package ollama

type OllamaOption func(*Ollama)

func WithUrl(url string) OllamaOption {
	return func(ollama *Ollama) {
		ollama.url = url
	}
}

func WithEmbeddingModel(embeddingModel string) OllamaOption {
	return func(ollama *Ollama) {
		ollama.embeddingModel = embeddingModel
	}
}

func WithChatModel(chatModel string) OllamaOption {
	return func(ollama *Ollama) {
		ollama.chatModel = chatModel
	}
}
