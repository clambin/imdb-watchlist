package unique_test

import (
	"cmp"
	"github.com/clambin/imdb-watchlist/pkg/unique"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"slices"
	"testing"
)

type record struct{ value int }

func getKey(r record) int { return r.value }

func TestUnique(t *testing.T) {
	var input []record
	var expected []record
	for i := range 10 {
		expected = append(expected, record{value: i})
		for range 1 + rand.Intn(5) {
			input = append(input, record{value: i})
		}
	}
	u := unique.UniqueFunc(input, getKey)
	slices.SortFunc(u, func(a, b record) int {
		return cmp.Compare(getKey(a), getKey(b))
	})
	assert.Equal(t, expected, u)
}

func BenchmarkUnique(b *testing.B) {
	const listSize = 15
	const duplicates = 3

	input := make([]record, listSize)
	for i := range listSize {
		input[i] = record{value: rand.Intn(listSize / duplicates)}
	}
	b.ResetTimer()
	for range b.N {
		unique.UniqueFunc(input, getKey)
	}
}
