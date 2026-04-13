// Package envfile provides utilities for parsing, writing, encrypting,
// and auditing .env files.
package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuditEvent represents a single recorded operation on an env file.
type AuditEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Environment string    `json:"environment"`
	Operation   string    `json:"operation"`
	Keys        []string  `json:"keys,omitempty"`
	Note        string    `json:"note,omitempty"`
}

// AuditLog holds an ordered list of audit events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

// LoadAuditLog reads an audit log from path. If the file does not exist
// an empty log is returned without error.
func LoadAuditLog(path string) (*AuditLog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read %s: %w", path, err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("audit: parse %s: %w", path, err)
	}
	return &log, nil
}

// Append adds a new event to the log and persists it to path.
func (l *AuditLog) Append(path string, event AuditEvent) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	l.Events = append(l.Events, event)
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("audit: write %s: %w", path, err)
	}
	return nil
}

// Record is a convenience helper that loads, appends, and saves in one call.
func Record(path, environment, operation string, keys []string, note string) error {
	log, err := LoadAuditLog(path)
	if err != nil {
		return err
	}
	return log.Append(path, AuditEvent{
		Timestamp:   time.Now().UTC(),
		Environment: environment,
		Operation:   operation,
		Keys:        keys,
		Note:        note,
	})
}
