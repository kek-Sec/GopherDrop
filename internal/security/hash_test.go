package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash_ValidLength(t *testing.T) {
	length := 16
	hash, err := GenerateHash(length)

	assert.NoError(t, err, "GenerateHash should not return an error")
	assert.Equal(t, length, len(hash), "Generated hash should have the correct length")
}

func TestGenerateHash_ZeroLength(t *testing.T) {
	length := 0
	hash, err := GenerateHash(length)

	assert.NoError(t, err, "GenerateHash should not return an error for zero length")
	assert.Equal(t, 0, len(hash), "Generated hash should be empty for zero length")
}
func TestGenerateHash_NegativeLength(t *testing.T) {
	length := -1
	hash, err := GenerateHash(length)

	assert.Error(t, err, "GenerateHash should return an error for negative length")
	assert.Empty(t, hash, "Generated hash should be empty for negative length")
}


func TestGenerateHash_UniqueHashes(t *testing.T) {
	length := 16
	hash1, err1 := GenerateHash(length)
	hash2, err2 := GenerateHash(length)

	assert.NoError(t, err1, "First GenerateHash call should not return an error")
	assert.NoError(t, err2, "Second GenerateHash call should not return an error")
	assert.NotEqual(t, hash1, hash2, "Generated hashes should be unique")
}
