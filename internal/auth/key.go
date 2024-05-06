package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateKey generates a Sonarr API Key
func GenerateKey() string {
	b := make([]byte, 16) // 16 bytes = 32 nibbles = 128 bits
	// in theory this could fail. On practice, it never happens.
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
