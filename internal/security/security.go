// Package security provides encryption and decryption functions.
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// PadKey ensures the key is 32 bytes.
func PadKey(k string) string {
	for len(k) < 32 {
		k += "0"
	}
	return k[:32]
}

// EncryptData encrypts data using AES-256 CFB mode.
func EncryptData(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	if len(data) > (1<<31 - 1 - aes.BlockSize) {
		return "", errors.New("data too large to encrypt")
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	copy(ciphertext[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts data previously encrypted with AES-256 CFB.
func DecryptData(enc string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
