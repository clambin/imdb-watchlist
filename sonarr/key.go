package sonarr

import (
	"math/rand"
	"time"
)

const (
	characters = `0123456789abcdef`
	size       = 32
)

// GenerateKey generates a Sonarr API Key
func GenerateKey() string {
	rand.Seed(time.Now().UnixNano())
	output := make([]byte, size)
	for i := 0; i < size; i++ {
		output[i] = characters[rand.Intn(len(characters))]
	}
	return string(output)
}
