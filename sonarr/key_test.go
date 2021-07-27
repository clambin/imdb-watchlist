package sonarr_test

import (
	"github.com/clambin/imdb-watchlist/sonarr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key1 := sonarr.GenerateKey()

	assert.Len(t, key1, 32)

	for i := 0; i < 100; i++ {
		key2 := sonarr.GenerateKey()
		if assert.NotEqual(t, key1, key2) == false {
			break
		}
	}
}
