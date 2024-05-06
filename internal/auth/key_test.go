package auth_test

import (
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	keys := make(map[string]struct{})

	for range 10000 {
		key := auth.GenerateKey()
		require.Len(t, key, 32)

		_, ok := keys[key]
		assert.False(t, ok)
		keys[key] = struct{}{}
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for range b.N {
		_ = auth.GenerateKey()
	}
}
