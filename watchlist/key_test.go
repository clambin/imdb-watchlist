package watchlist_test

import (
	"github.com/clambin/imdb-watchlist/watchlist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	keys := make(map[string]struct{})

	for i := 0; i < 1000; i++ {
		key := watchlist.GenerateKey()
		assert.Len(t, key, 32)

		_, ok := keys[key]
		assert.False(t, ok)
		keys[key] = struct{}{}
	}
}
