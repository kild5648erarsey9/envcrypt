package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/internal/keystore"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "keys.json")
}

func TestLoadEmptyStore(t *testing.T) {
	ks, err := keystore.Load(tempPath(t))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ks.List()) != 0 {
		t.Fatalf("expected empty store")
	}
}

func TestSetAndGet(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))
	ks.Set("production", "deadbeefdeadbeef")

	e, ok := ks.Get("production")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Key != "deadbeefdeadbeef" {
		t.Errorf("unexpected key: %s", e.Key)
	}
	if e.RotatedAt != nil {
		t.Error("new entry should not have RotatedAt set")
	}
}

func TestRotateUpdatesKey(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))
	ks.Set("staging", "aabbccdd")
	ks.Set("staging", "11223344")

	e, _ := ks.Get("staging")
	if e.Key != "11223344" {
		t.Errorf("expected rotated key, got %s", e.Key)
	}
	if e.RotatedAt == nil {
		t.Error("RotatedAt should be set after rotation")
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tempPath(t)
	ks, _ := keystore.Load(path)
	ks.Set("dev", "cafebabe")
	if err := ks.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	ks2, err := keystore.Load(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	e, ok := ks2.Get("dev")
	if !ok || e.Key != "cafebabe" {
		t.Errorf("expected reloaded key, got %+v", e)
	}
}

func TestDelete(t *testing.T) {
	ks, _ := keystore.Load(tempPath(t))
	ks.Set("test", "key123")

	if !ks.Delete("test") {
		t.Error("expected Delete to return true")
	}
	if _, ok := ks.Get("test"); ok {
		t.Error("entry should not exist after delete")
	}
	if ks.Delete("test") {
		t.Error("expected Delete to return false for missing entry")
	}
}

func TestSaveCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "keys.json")
	ks, _ := keystore.Load(path)
	ks.Set("env", "hexkey")
	if err := ks.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
