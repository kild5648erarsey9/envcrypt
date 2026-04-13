package envfile

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds all issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d issue(s): %s", len(e.Issues), strings.Join(e.Issues, "; "))
}

// validKeyRe matches POSIX-style env var names: uppercase letters, digits, underscores,
// must not start with a digit.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks a slice of Entry values for common problems:
//   - empty keys
//   - keys that do not match the valid identifier pattern
//   - duplicate keys
//
// Returns a *ValidationError if any issues are found, nil otherwise.
func Validate(entries []Entry) error {
	var issues []string
	seen := make(map[string]int) // key -> first line index (1-based)

	for i, e := range entries {
		lineNum := i + 1

		if e.Key == "" {
			issues = append(issues, fmt.Sprintf("line %d: empty key", lineNum))
			continue
		}

		if !validKeyRe.MatchString(e.Key) {
			issues = append(issues, fmt.Sprintf("line %d: invalid key %q (must match [A-Za-z_][A-Za-z0-9_]*)", lineNum, e.Key))
		}

		if first, dup := seen[e.Key]; dup {
			issues = append(issues, fmt.Sprintf("line %d: duplicate key %q (first seen on line %d)", lineNum, e.Key, first))
		} else {
			seen[e.Key] = lineNum
		}
	}

	if len(issues) == 0 {
		return nil
	}
	return &ValidationError{Issues: issues}
}

// ValidateFile parses the file at path and then validates its entries.
func ValidateFile(path string) ([]Entry, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	if err := Validate(entries); err != nil {
		return entries, err
	}
	return entries, nil
}

// ErrInvalidKey is a sentinel used when a single key is checked.
var ErrInvalidKey = errors.New("invalid environment variable key")

// ValidateKey returns ErrInvalidKey if the provided key string is not a valid
// environment variable identifier.
func ValidateKey(key string) error {
	if !validKeyRe.MatchString(key) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	return nil
}
