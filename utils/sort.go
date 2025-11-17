package utils

import "sort"

func SortedValuesByKey[V any](m map[string]V) []V {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	out := make([]V, 0, len(keys))
	for _, k := range keys {
		out = append(out, m[k])
	}

	return out
}
