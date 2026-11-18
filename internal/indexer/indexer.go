package indexer

import "iter"

type Indexer struct {
	chunkSize    int
	chunkOverlap int
}

type SearchResult struct {
	Text   string
	Start  int
	Length int
}

func NewIndexer(opts ...IndexerOption) *Indexer {
	indexer := &Indexer{
		chunkSize:    500,
		chunkOverlap: 100,
	}
	for _, opt := range opts {
		opt(indexer)
	}
	return indexer
}

func (indexer *Indexer) Split(text string) iter.Seq[SearchResult] {
	return func(yield func(SearchResult) bool) {
		step := indexer.chunkSize - indexer.chunkOverlap
		for left := 0; left < len(text)-indexer.chunkSize; left += step {
			searchResult := SearchResult{
				Text:   text[left : left+indexer.chunkSize],
				Start:  left,
				Length: indexer.chunkSize,
			}
			if !yield(searchResult) {
				return
			}
		}
	}
}
