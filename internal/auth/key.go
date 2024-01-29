package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateKey generates a Sonarr API Key
func GenerateKey() (string, error) {
	var key string
	b := make([]byte, 16) // 16 bytes = 32 nibbles = 128 bits
	_, err := rand.Read(b)
	if err == nil {
		key = hex.EncodeToString(b)
	}
	return key, nil
}
