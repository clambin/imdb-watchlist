package auth_test

import (
	"github.com/clambin/go-common/set"
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	keys := set.New[string]()

	for range 10000 {
		key, err := auth.GenerateKey()
		require.NoError(t, err)
		require.Len(t, key, 32)

		assert.False(t, keys.Contains(key))
		keys.Add(key)
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for range b.N {
		_, _ = auth.GenerateKey()
	}
}
