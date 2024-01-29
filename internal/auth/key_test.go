package auth_test

import (
	"github.com/clambin/imdb-watchlist/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	keys := make(map[string]struct{})

	for i := 0; i < 1e4; i++ {
		key, err := auth.GenerateKey()
		require.NoError(t, err)
		require.Len(t, key, 32)

		_, ok := keys[key]
		assert.False(t, ok)
		keys[key] = struct{}{}
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = auth.GenerateKey()
	}
}
