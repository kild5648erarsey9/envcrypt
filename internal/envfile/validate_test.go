package envfile

import (
	"os"
	"strings"
	"testing"
)

func TestValidateClean(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "_UNDERSCORE", Value: "ok"},
		{Key: "MixedCase123", Value: "fine"},
	}
	if err := Validate(entries); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateEmptyKey(t *testing.T) {
	entries := []Entry{{Key: "", Value: "oops"}}
	err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 || !strings.Contains(ve.Issues[0], "empty key") {
		t.Errorf("unexpected issues: %v", ve.Issues)
	}
}

func TestValidateInvalidKey(t *testing.T) {
	entries := []Entry{{Key: "1INVALID", Value: "v"}}
	err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("error message should mention 'invalid key', got: %v", err)
	}
}

func TestValidateDuplicateKeys(t *testing.T) {
	entries := []Entry{
		{Key: "DUP", Value: "first"},
		{Key: "OTHER", Value: "x"},
		{Key: "DUP", Value: "second"},
	}
	err := Validate(entries)
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	ve := err.(*ValidationError)
	if len(ve.Issues) != 1 || !strings.Contains(ve.Issues[0], "duplicate key") {
		t.Errorf("unexpected issues: %v", ve.Issues)
	}
}

func TestValidateMultipleIssues(t *testing.T) {
	entries := []Entry{
		{Key: "GOOD", Value: "ok"},
		{Key: "bad-key", Value: "v"},
		{Key: "GOOD", Value: "dup"},
	}
	err := Validate(entries)
	if err == nil {
		t.Fatal("expected error")
	}
	ve := err.(*ValidationError)
	if len(ve.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(ve.Issues), ve.Issues)
	}
}

func TestValidateFile(t *testing.T) {
	f, err := os.CreateTemp("", "validate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("VALID_KEY=hello\nANOTHER=world\n")
	f.Close()

	entries, err := ValidateFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestValidateKey(t *testing.T) {
	if err := ValidateKey("GOOD_KEY"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if err := ValidateKey("bad-key"); err == nil {
		t.Error("expected error for bad-key")
	}
}
