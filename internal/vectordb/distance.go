package vectordb

import (
	"math"
)

func Euclidean(a, b []float64) float64 {
	var sum float64

	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

func SquaredEuclidean(a, b []float64) float64 {
	var sum float64

	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sum
}

func Manhattan(a, b []float64) float64 {
	var sum float64

	for i := range a {
		sum += math.Abs(a[i] - b[i])
	}

	return sum
}

func DotProduct(a, b []float64) float64 {
	var sum float64

	for i := range a {
		sum += a[i] * b[i]
	}

	return sum
}

func CosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64

	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func CosineDistance(a, b []float64) float64 {
	return 1 - CosineSimilarity(a, b)
}
