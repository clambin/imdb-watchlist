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

	for i := 0; i < 1e4; i++ {
		key, err := auth.GenerateKey()
		require.NoError(t, err)
		require.Len(t, key, 32)

		assert.False(t, keys.Contains(key))
		keys.Add(key)
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = auth.GenerateKey()
	}
}
