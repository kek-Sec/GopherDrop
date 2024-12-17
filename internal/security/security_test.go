package security

import (
	"testing"
)

func TestPadKey(t *testing.T) {
	key := PadKey("short")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}
}

func TestEncryptDecryptData(t *testing.T) {
	key := []byte(PadKey("password"))
	original := []byte("hello world")
	enc, err := EncryptData(original, key)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}
	dec, err := DecryptData(enc, key)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}
	if string(dec) != string(original) {
		t.Fatalf("expected '%s' got '%s'", original, dec)
	}
}

func TestGenerateHash(t *testing.T) {
	hash, err := GenerateHash(16)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}
	if len(hash) != 16 {
		t.Fatalf("expected length 16, got %d", len(hash))
	}
}
