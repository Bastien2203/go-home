package utils

func Map[T any, K any](list []T, modifier func(T) K) []K {
	newList := make([]K, len(list))
	for i, v := range list {
		newList[i] = modifier(v)
	}
	return newList
}

func MapEntries[K comparable, V any, T any](entries map[K]V, modifier func(K, V) T) []T {
	list := make([]T, 0, len(entries))
	for k, v := range entries {
		list = append(list, modifier(k, v))
	}
	return list
}
