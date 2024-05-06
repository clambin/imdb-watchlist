package unique

import (
	"cmp"
	"slices"
)

// UniqueFunc returns the unique values from the slice, using the getKey function to determine the uniqueness of the element in the slice.
// The returned slice is not guaranteed to be in order.
func UniqueFunc[S ~[]E, E any, K cmp.Ordered](input S, getKey func(entry E) K) S {
	slices.SortFunc(input, func(a, b E) int {
		return cmp.Compare(getKey(a), getKey(b))
	})
	return slices.CompactFunc(input, func(v E, v2 E) bool {
		return getKey(v) == getKey(v2)
	})
}
