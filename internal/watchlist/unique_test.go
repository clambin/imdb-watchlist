package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/internal/watchlist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnique(t *testing.T) {
	output := watchlist.Unique([]int{5, 4, 3, 2, 1, 5}, func(v int) int {
		return v
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, output)

	type structured struct {
		id string
	}

	input := []structured{{id: "3"}, {id: "2"}, {id: "1"}, {id: "0"}, {id: "3"}, {id: "2"}}
	want := []structured{{id: "0"}, {id: "1"}, {id: "2"}, {id: "3"}}

	assert.Equal(t, want, watchlist.Unique(input, func(v structured) string {
		return v.id
	}))
}

func BenchmarkUnique(b *testing.B) {
	const count = 10000
	bigInput := make([]int, count)
	for i := range count {
		bigInput[i] = i
	}

	b.ResetTimer()
	for range b.N {
		_ = watchlist.Unique(bigInput, func(v int) int { return v })
	}
}
