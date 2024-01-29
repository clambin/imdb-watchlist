package watchlist

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateKey generates a Sonarr API Key
func GenerateKey() string {
	b := make([]byte, 16) // 16 bytes = 32 nibbles = 128 bits
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
