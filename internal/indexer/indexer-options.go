package indexer

type IndexerOption func(*Indexer)

func WithChunkSize(chunkSize int) IndexerOption {
	return func(indexer *Indexer) {
		indexer.chunkSize = chunkSize
	}
}

func WithChunkOverlap(chunkOverlap int) IndexerOption {
	return func(indexer *Indexer) {
		indexer.chunkOverlap = chunkOverlap
	}
}
