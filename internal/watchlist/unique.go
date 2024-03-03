package watchlist

import (
	"cmp"
	"slices"
)

func Unique[V any, K cmp.Ordered](input []V, getKey func(V) K) []V {
	slices.SortFunc(input, func(a, b V) int {
		return cmp.Compare(getKey(a), getKey(b))
	})
	var last K
	entries := make([]V, 0, len(input))
	for _, e := range input {
		key := getKey(e)
		if key == last {
			continue
		}
		entries = append(entries, e)
		last = key
	}
	return entries
}
