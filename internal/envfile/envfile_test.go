package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestParseBasic(t *testing.T) {
	p := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "DB_HOST" || ef.Entries[0].Value != "localhost" {
		t.Errorf("unexpected first entry: %+v", ef.Entries[0])
	}
}

func TestParseSkipsCommentsAndBlanks(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nFOO=bar\n")
	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
}

func TestParseMalformedLine(t *testing.T) {
	p := writeTempEnv(t, "NOEQUALSSIGN\n")
	_, err := Parse(p)
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestWriteRoundTrip(t *testing.T) {
	original := &EnvFile{
		Entries: []Entry{
			{Key: "APP_ENV", Value: "production"},
			{Key: "SECRET", Value: "s3cr3t"},
		},
	}
	dir := t.TempDir()
	p := filepath.Join(dir, "out.env")
	if err := Write(p, original); err != nil {
		t.Fatalf("Write: %v", err)
	}
	loaded, err := Parse(p)
	if err != nil {
		t.Fatalf("Parse after Write: %v", err)
	}
	if len(loaded.Entries) != len(original.Entries) {
		t.Fatalf("entry count mismatch: want %d got %d", len(original.Entries), len(loaded.Entries))
	}
	for i, e := range original.Entries {
		if loaded.Entries[i] != e {
			t.Errorf("entry %d mismatch: want %+v got %+v", i, e, loaded.Entries[i])
		}
	}
}
