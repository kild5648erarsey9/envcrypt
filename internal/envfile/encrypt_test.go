package envfile

import (
	"testing"

	"envcrypt/internal/crypto"
)

func newTestEnvFile() *EnvFile {
	return &EnvFile{
		Entries: []Entry{
			{Key: "DB_PASSWORD", Value: "supersecret"},
			{Key: "API_KEY", Value: "abc123xyz"},
		},
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	ef := newTestEnvFile()
	originalValues := make([]string, len(ef.Entries))
	for i, e := range ef.Entries {
		originalValues[i] = e.Value
	}

	if err := EncryptValues(ef, key); err != nil {
		t.Fatalf("EncryptValues: %v", err)
	}
	// Values should now differ from originals.
	for i, e := range ef.Entries {
		if e.Value == originalValues[i] {
			t.Errorf("entry %q: value not encrypted", e.Key)
		}
	}

	if err := DecryptValues(ef, key); err != nil {
		t.Fatalf("DecryptValues: %v", err)
	}
	for i, e := range ef.Entries {
		if e.Value != originalValues[i] {
			t.Errorf("entry %q: want %q got %q", e.Key, originalValues[i], e.Value)
		}
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key, _ := crypto.GenerateKey()
	wrongKey, _ := crypto.GenerateKey()

	ef := newTestEnvFile()
	if err := EncryptValues(ef, key); err != nil {
		t.Fatalf("EncryptValues: %v", err)
	}
	if err := DecryptValues(ef, wrongKey); err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	key, _ := crypto.GenerateKey()
	ef := &EnvFile{
		Entries: []Entry{{Key: "X", Value: "!!!notbase64!!!"}},
	}
	if err := DecryptValues(ef, key); err == nil {
		t.Fatal("expected error for invalid base64 value")
	}
}
