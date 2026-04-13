package keystore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// KeyEntry represents a stored encryption key for a specific environment.
type KeyEntry struct {
	Environment string    `json:"environment"`
	Key         string    `json:"key"` // hex-encoded AES key
	CreatedAt   time.Time `json:"created_at"`
	RotatedAt   *time.Time `json:"rotated_at,omitempty"`
}

// KeyStore manages per-environment encryption keys.
type KeyStore struct {
	path    string
	entries map[string]KeyEntry
}

// Load reads the keystore from disk, or returns an empty store if not found.
func Load(path string) (*KeyStore, error) {
	ks := &KeyStore{
		path:    path,
		entries: make(map[string]KeyEntry),
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ks, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []KeyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		ks.entries[e.Environment] = e
	}
	return ks, nil
}

// Save persists the keystore to disk.
func (ks *KeyStore) Save() error {
	if err := os.MkdirAll(filepath.Dir(ks.path), 0700); err != nil {
		return err
	}
	entries := make([]KeyEntry, 0, len(ks.entries))
	for _, e := range ks.entries {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ks.path, data, 0600)
}

// Set stores or updates the key for an environment.
func (ks *KeyStore) Set(env, hexKey string) {
	now := time.Now().UTC()
	if existing, ok := ks.entries[env]; ok {
		existing.Key = hexKey
		existing.RotatedAt = &now
		ks.entries[env] = existing
		return
	}
	ks.entries[env] = KeyEntry{
		Environment: env,
		Key:         hexKey,
		CreatedAt:   now,
	}
}

// Get retrieves the key entry for an environment.
func (ks *KeyStore) Get(env string) (KeyEntry, bool) {
	e, ok := ks.entries[env]
	return e, ok
}

// List returns all stored key entries.
func (ks *KeyStore) List() []KeyEntry {
	result := make([]KeyEntry, 0, len(ks.entries))
	for _, e := range ks.entries {
		result = append(result, e)
	}
	return result
}

// Delete removes the key entry for an environment.
func (ks *KeyStore) Delete(env string) bool {
	if _, ok := ks.entries[env]; !ok {
		return false
	}
	delete(ks.entries, env)
	return true
}
