package oracle

type OracleOption func(*Oracle)

func WithTopKEmbeddings(topk int) OracleOption {
	return func(oracle *Oracle) {
		oracle.topKEmbeddings = topk
	}
}
