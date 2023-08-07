package slice

func Map[T, V any](collection []T, fn func(T) V) []V {
	result := make([]V, len(collection))
	for i, t := range collection {
		result[i] = fn(t)
	}
	return result
}

func Reduce[T, V any](collection []T, accumulator func(V, T) V, initValue V) V {
	var result = initValue
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}
