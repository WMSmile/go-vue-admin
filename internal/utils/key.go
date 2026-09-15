package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAPIKey returns a cryptographically random API key, prefixed with "gva_".
func GenerateAPIKey() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "gva_" + hex.EncodeToString(b), nil
}
