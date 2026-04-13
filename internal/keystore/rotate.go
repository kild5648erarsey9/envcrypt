package keystore

import (
	"encoding/hex"
	"fmt"

	"github.com/user/envcrypt/internal/crypto"
)

// RotateResult holds the old and new key after a rotation.
type RotateResult struct {
	Environment string
	OldKey      string
	NewKey      string
}

// Rotate generates a new AES key for the given environment, stores it in the
// keystore, and returns the old and new keys so the caller can re-encrypt data.
func Rotate(ks *KeyStore, env string) (RotateResult, error) {
	result := RotateResult{Environment: env}

	if existing, ok := ks.Get(env); ok {
		result.OldKey = existing.Key
	}

	newKeyBytes, err := crypto.GenerateKey()
	if err != nil {
		return result, fmt.Errorf("rotate: generate key: %w", err)
	}

	newKeyHex := hex.EncodeToString(newKeyBytes)
	ks.Set(env, newKeyHex)
	result.NewKey = newKeyHex
	return result, nil
}

// GenerateAndStore creates a new key for the environment only if one does not
// already exist. Returns an error if the environment already has a key.
func GenerateAndStore(ks *KeyStore, env string) (string, error) {
	if _, ok := ks.Get(env); ok {
		return "", fmt.Errorf("key for environment %q already exists; use rotate to replace it", env)
	}

	keyBytes, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}

	hexKey := hex.EncodeToString(keyBytes)
	ks.Set(env, hexKey)
	return hexKey, nil
}
