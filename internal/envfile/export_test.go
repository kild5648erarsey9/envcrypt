package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "s3cr3t"},
	}
}

func TestExportDotenv(t *testing.T) {
	out := exportDotenv(testEntries())
	if !strings.Contains(out, "APP_NAME=envcrypt") {
		t.Errorf("expected APP_NAME=envcrypt in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASS=s3cr3t") {
		t.Errorf("expected DB_PASS in output")
	}
}

func TestExportShell(t *testing.T) {
	out := exportShell(testEntries())
	if !strings.Contains(out, "export APP_NAME=") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExportJSON(t *testing.T) {
	out, err := exportJSON(testEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["APP_NAME"] != "envcrypt" {
		t.Errorf("expected APP_NAME=envcrypt, got %s", m["APP_NAME"])
	}
}

func TestExportPrefixFilter(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.env")
	err := Export(testEntries(), ExportOptions{Format: FormatDotenv, Prefix: "DB_"}, tmp)
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if strings.Contains(string(data), "APP_NAME") {
		t.Error("APP_NAME should have been filtered out")
	}
	if !strings.Contains(string(data), "DB_HOST") {
		t.Error("DB_HOST should be present")
	}
}

func TestExportRedact(t *testing.T) {
	out := exportDotenv(func() []Entry {
		e := testEntries()
		for i := range e {
			e[i].Value = "***"
		}
		return e
	}())
	if strings.Contains(out, "s3cr3t") {
		t.Error("redacted export must not contain real values")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected redaction placeholder")
	}
}

func TestExportToFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "export.json")
	err := Export(testEntries(), ExportOptions{Format: FormatJSON}, tmp)
	if err != nil {
		t.Fatalf("Export to file: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid JSON file: %v", err)
	}
}
