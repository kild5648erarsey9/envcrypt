package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempAuditPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestLoadAuditLogMissing(t *testing.T) {
	log, err := LoadAuditLog("/nonexistent/audit.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(log.Events) != 0 {
		t.Fatalf("expected empty log, got %d events", len(log.Events))
	}
}

func TestAppendAndReload(t *testing.T) {
	path := tempAuditPath(t)
	log, _ := LoadAuditLog(path)

	event := AuditEvent{
		Timestamp:   time.Now().UTC(),
		Environment: "production",
		Operation:   "encrypt",
		Keys:        []string{"DB_PASSWORD", "API_KEY"},
	}
	if err := log.Append(path, event); err != nil {
		t.Fatalf("Append: %v", err)
	}

	reloaded, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(reloaded.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(reloaded.Events))
	}
	got := reloaded.Events[0]
	if got.Environment != "production" {
		t.Errorf("environment: want production, got %s", got.Environment)
	}
	if got.Operation != "encrypt" {
		t.Errorf("operation: want encrypt, got %s", got.Operation)
	}
	if len(got.Keys) != 2 {
		t.Errorf("keys: want 2, got %d", len(got.Keys))
	}
}

func TestAppendSetsTimestampWhenZero(t *testing.T) {
	path := tempAuditPath(t)
	log, _ := LoadAuditLog(path)

	before := time.Now().UTC()
	_ = log.Append(path, AuditEvent{Environment: "staging", Operation: "rotate"})
	after := time.Now().UTC()

	if log.Events[0].Timestamp.Before(before) || log.Events[0].Timestamp.After(after) {
		t.Error("timestamp not set correctly for zero-value event")
	}
}

func TestRecordMultipleEvents(t *testing.T) {
	path := tempAuditPath(t)

	ops := []string{"init", "encrypt", "rotate"}
	for _, op := range ops {
		if err := Record(path, "dev", op, nil, ""); err != nil {
			t.Fatalf("Record %s: %v", op, err)
		}
	}

	log, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(log.Events))
	}
	for i, op := range ops {
		if log.Events[i].Operation != op {
			t.Errorf("event[%d] operation: want %s, got %s", i, op, log.Events[i].Operation)
		}
	}
}

func TestLoadAuditLogInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o600)
	_, err := LoadAuditLog(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
