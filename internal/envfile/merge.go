// Package envfile provides utilities for parsing, encrypting, and managing .env files.
package envfile

import "fmt"

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the value from the base (current) file on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the value from the incoming file on conflict.
	MergeStrategyTheirs
	// MergeStrategyError returns an error when a conflict is detected.
	MergeStrategyError
)

// MergeResult holds the merged key-value pairs and metadata about the operation.
type MergeResult struct {
	Entries    []Entry
	Added      []string
	Overridden []string
	Skipped    []string
}

// Merge combines base and incoming env entries according to the given strategy.
// Keys present only in incoming are always added. Keys present in both are
// resolved using the strategy.
func Merge(base, incoming []Entry, strategy MergeStrategy) (*MergeResult, error) {
	result := &MergeResult{}

	// Build a map from base for quick lookup.
	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	// Start with all base entries.
	merged := make(map[string]string, len(base))
	for _, e := range base {
		merged[e.Key] = e.Value
	}

	for _, e := range incoming {
		if _, exists := baseMap[e.Key]; exists {
			if baseMap[e.Key] == e.Value {
				// Same value — no conflict.
				continue
			}
			switch strategy {
			case MergeStrategyOurs:
				result.Skipped = append(result.Skipped, e.Key)
			case MergeStrategyTheirs:
				merged[e.Key] = e.Value
				result.Overridden = append(result.Overridden, e.Key)
			case MergeStrategyError:
				return nil, fmt.Errorf("merge conflict on key %q", e.Key)
			}
		} else {
			merged[e.Key] = e.Value
			result.Added = append(result.Added, e.Key)
		}
	}

	// Preserve order: base first, then newly added keys in incoming order.
	addedSet := make(map[string]bool, len(result.Added))
	for _, k := range result.Added {
		addedSet[k] = true
	}

	for _, e := range base {
		result.Entries = append(result.Entries, Entry{Key: e.Key, Value: merged[e.Key]})
	}
	for _, e := range incoming {
		if addedSet[e.Key] {
			result.Entries = append(result.Entries, Entry{Key: e.Key, Value: merged[e.Key]})
		}
	}

	return result, nil
}
