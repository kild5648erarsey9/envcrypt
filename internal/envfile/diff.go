package envfile

import "sort"

// DiffResult holds the changes between two env file snapshots.
type DiffResult struct {
	Added   map[string]string // keys present in new but not old
	Removed map[string]string // keys present in old but not new
	Changed map[string][2]string // keys whose values changed: [old, new]
}

// Diff compares two parsed env maps and returns what changed.
// Both arguments should be the result of Parse or DecryptValues.
func Diff(oldEnv, newEnv map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, newVal := range newEnv {
		oldVal, exists := oldEnv[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = [2]string{oldVal, newVal}
		}
	}

	for k, oldVal := range oldEnv {
		if _, exists := newEnv[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// IsEmpty reports whether the DiffResult contains no changes.
func (d DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// SortedAdded returns the added keys in sorted order.
func (d DiffResult) SortedAdded() []string {
	return sortedKeys(d.Added)
}

// SortedRemoved returns the removed keys in sorted order.
func (d DiffResult) SortedRemoved() []string {
	return sortedKeys(d.Removed)
}

// SortedChanged returns the changed keys in sorted order.
func (d DiffResult) SortedChanged() []string {
	keys := make([]string, 0, len(d.Changed))
	for k := range d.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
