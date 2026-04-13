package envfile

import (
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMergeNoConflict(t *testing.T) {
	base := entries("A", "1", "B", "2")
	incoming := entries("C", "3")
	res, err := Merge(base, incoming, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Errorf("expected C to be added, got %v", res.Added)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
}

func TestMergeStrategyOursKeepsBase(t *testing.T) {
	base := entries("A", "original")
	incoming := entries("A", "override")
	res, err := Merge(base, incoming, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	if res.Entries[0].Value != "original" {
		t.Errorf("expected original value, got %q", res.Entries[0].Value)
	}
}

func TestMergeStrategyTheirsOverrides(t *testing.T) {
	base := entries("A", "original")
	incoming := entries("A", "override")
	res, err := Merge(base, incoming, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overridden) != 1 || res.Overridden[0] != "A" {
		t.Errorf("expected A to be overridden, got %v", res.Overridden)
	}
	if res.Entries[0].Value != "override" {
		t.Errorf("expected overridden value, got %q", res.Entries[0].Value)
	}
}

func TestMergeStrategyErrorOnConflict(t *testing.T) {
	base := entries("A", "1")
	incoming := entries("A", "2")
	_, err := Merge(base, incoming, MergeStrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMergeSameValueNoConflict(t *testing.T) {
	base := entries("A", "same")
	incoming := entries("A", "same")
	res, err := Merge(base, incoming, MergeStrategyError)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
	if len(res.Overridden) != 0 || len(res.Skipped) != 0 {
		t.Errorf("expected no conflict metadata for identical values")
	}
}

func TestMergePreservesOrder(t *testing.T) {
	base := entries("B", "2", "A", "1")
	incoming := entries("C", "3", "D", "4")
	res, err := Merge(base, incoming, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	expected := []string{"B", "A", "C", "D"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}
