package vectordb

type MemoryOption func(*Memory)

func MemoryWithDistanceFunc(distanceMetric DistanceMetric) MemoryOption {
	return func(memory *Memory) {
		memory.distanceMetric = distanceMetric
	}
}
