package keystore_test

import (
	"testing"

	"github.com/user/envcrypt/internal/keystore"
)

func TestGenerateAndStore(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))

	hexKey, err := keystore.GenerateAndStore(ks, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hexKey) != 64 { // 32 bytes -> 64 hex chars
		t.Errorf("expected 64 hex chars, got %d", len(hexKey))
	}

	e, ok := ks.Get("production")
	if !ok {
		t.Fatal("key not stored")
	}
	if e.Key != hexKey {
		t.Errorf("stored key mismatch")
	}
}

func TestGenerateAndStoreDuplicateErrors(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))
	if _, err := keystore.GenerateAndStore(ks, "staging"); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if _, err := keystore.GenerateAndStore(ks, "staging"); err == nil {
		t.Fatal("expected error on duplicate, got nil")
	}
}

func TestRotateReplacesKey(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))

	original, _ := keystore.GenerateAndStore(ks, "dev")

	result, err := keystore.Rotate(ks, "dev")
	if err != nil {
		t.Fatalf("rotate failed: %v", err)
	}
	if result.OldKey != original {
		t.Errorf("expected old key %s, got %s", original, result.OldKey)
	}
	if result.NewKey == original {
		t.Error("new key should differ from old key")
	}
	if len(result.NewKey) != 64 {
		t.Errorf("expected 64 hex chars for new key, got %d", len(result.NewKey))
	}

	e, _ := ks.Get("dev")
	if e.Key != result.NewKey {
		t.Error("keystore not updated with new key")
	}
}

func TestRotateNewEnvironment(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))

	result, err := keystore.Rotate(ks, "newenv")
	if err != nil {
		t.Fatalf("rotate on new env failed: %v", err)
	}
	if result.OldKey != "" {
		t.Errorf("expected empty old key for new env, got %s", result.OldKey)
	}
	if result.NewKey == "" {
		t.Error("expected non-empty new key")
	}
}
