package security

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateHash creates a secure random hash of the specified length.
func GenerateHash(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
