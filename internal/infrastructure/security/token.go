package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateRandomToken generates a cryptographically secure random hex string of n bytes (2n hex chars)
func GenerateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken computes a SHA-256 hex digest of a token for safe database storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
