package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateKey generates a Sonarr API Key
func GenerateKey() (string, error) {
	b := make([]byte, 16) // 16 bytes = 32 nibbles = 128 bits
	_, err := rand.Read(b)
	var key string
	if err == nil {
		key = hex.EncodeToString(b)
	}
	return key, err
}
