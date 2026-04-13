package envfile

import (
	"encoding/base64"
	"fmt"

	"envcrypt/internal/crypto"
)

// EncryptValues encrypts every entry value in ef using the provided AES key.
// Each value is replaced with a base64-encoded ciphertext so the result is
// safe to store as a plain text .env file.
func EncryptValues(ef *EnvFile, key []byte) error {
	for i, e := range ef.Entries {
		ciphertext, err := crypto.Encrypt(key, []byte(e.Value))
		if err != nil {
			return fmt.Errorf("envfile: encrypt key %q: %w", e.Key, err)
		}
		ef.Entries[i].Value = base64.StdEncoding.EncodeToString(ciphertext)
	}
	return nil
}

// DecryptValues decrypts every entry value in ef using the provided AES key.
// Values are expected to be base64-encoded ciphertexts produced by EncryptValues.
func DecryptValues(ef *EnvFile, key []byte) error {
	for i, e := range ef.Entries {
		ciphertext, err := base64.StdEncoding.DecodeString(e.Value)
		if err != nil {
			return fmt.Errorf("envfile: base64 decode key %q: %w", e.Key, err)
		}
		plaintext, err := crypto.Decrypt(key, ciphertext)
		if err != nil {
			return fmt.Errorf("envfile: decrypt key %q: %w", e.Key, err)
		}
		ef.Entries[i].Value = string(plaintext)
	}
	return nil
}
