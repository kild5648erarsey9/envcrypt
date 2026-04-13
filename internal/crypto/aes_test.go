package crypto_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envcrypt/internal/crypto"
)

func TestGenerateKey(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected key length 32, got %d", len(key))
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	plaintext := []byte("DATABASE_URL=postgres://user:pass@localhost/mydb")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()

	plaintext := []byte("SECRET=supersecret")
	ciphertext, err := crypto.Encrypt(key1, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	_, err = crypto.Decrypt(key2, ciphertext)
	if err == nil {
		t.Error("expected error when decrypting with wrong key, got nil")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	key, _ := crypto.GenerateKey()
	_, err := crypto.Decrypt(key, []byte("short"))
	if err == nil {
		t.Error("expected error for short ciphertext, got nil")
	}
}

func TestEncryptProducesUniqueNonces(t *testing.T) {
	key, _ := crypto.GenerateKey()
	plaintext := []byte("SAME_CONTENT=true")

	ct1, _ := crypto.Encrypt(key, plaintext)
	ct2, _ := crypto.Encrypt(key, plaintext)

	if bytes.Equal(ct1, ct2) {
		t.Error("two encryptions of the same plaintext should not produce identical ciphertext")
	}
}
