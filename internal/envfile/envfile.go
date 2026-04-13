// Package envfile provides functionality to parse, encrypt, and decrypt
// .env files using per-environment AES keys managed by the keystore.
package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// EnvFile holds the parsed entries of a .env file in order.
type EnvFile struct {
	Entries []Entry
}

// Parse reads a .env file from the given path and returns an EnvFile.
// Lines beginning with '#' and empty lines are skipped.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	var ef EnvFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("envfile: malformed line: %q", line)
		}
		ef.Entries = append(ef.Entries, Entry{
			Key:   strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scan: %w", err)
	}
	return &ef, nil
}

// Write serialises the EnvFile entries to the given path, one KEY=VALUE per line.
func Write(path string, ef *EnvFile) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envfile: create %q: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range ef.Entries {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value); err != nil {
			return fmt.Errorf("envfile: write: %w", err)
		}
	}
	return w.Flush()
}
