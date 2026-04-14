package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for exported env data.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatShell  ExportFormat = "shell"
)

// ExportOptions configures the export behaviour.
type ExportOptions struct {
	Format  ExportFormat
	Redact  bool   // replace values with "***"
	Prefix  string // optional prefix filter (e.g. "DB_")
}

// Export writes env entries to w in the requested format.
// Only entries whose keys begin with opts.Prefix are included (all if empty).
func Export(entries []Entry, opts ExportOptions, path string) error {
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if opts.Redact {
			e.Value = "***"
		}
		filtered = append(filtered, e)
	}

	var out string
	var err error

	switch opts.Format {
	case FormatJSON:
		out, err = exportJSON(filtered)
	case FormatShell:
		out = exportShell(filtered)
	default:
		out = exportDotenv(filtered)
	}
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if path == "" || path == "-" {
		fmt.Print(out)
		return nil
	}
	return os.WriteFile(path, []byte(out), 0o600)
}

func exportDotenv(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func exportShell(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
	}
	return sb.String()
}

func exportJSON(entries []Entry) (string, error) {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
