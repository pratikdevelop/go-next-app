package helpers

import (
	"crypto/rand"
	"fmt"
)

// charset defines the characters used for generating short codes.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode creates a random short code of the specified length.
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	// Generate random bytes
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Map random bytes to charset
	result := make([]byte, length)
	for i := range b {
		result[i] = charset[int(b[i])%len(charset)]
	}

	return string(result), nil
}
