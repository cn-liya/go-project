package types

func SliceToMap[T comparable](arr []T) map[T]struct{} {
	m := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	return m
}

func SliceToCountMap[T comparable](arr []T) map[T]int {
	m := make(map[T]int, len(arr))
	for _, v := range arr {
		m[v]++
	}
	return m
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
