// Package-level documentation addition for the export feature.
// This file keeps doc.go clean while describing the export sub-feature.
package envfile

// Entry is the canonical key/value pair used throughout the envfile package.
// It is produced by Parse and consumed by Write, Export, Validate, Merge, and Diff.
//
// Export formats
//
// Three export formats are supported:
//
//	FormatDotenv  – plain KEY=VALUE lines (default, compatible with most tools)
//	FormatShell   – export KEY="VALUE" lines suitable for eval in bash/zsh
//	FormatJSON    – a JSON object mapping keys to values
//
// Filtering and redaction
//
// ExportOptions.Prefix limits output to entries whose keys start with the
// given string.  ExportOptions.Redact replaces every value with "***" so the
// resulting file can be committed or shared without leaking secrets.
//
// Output destination
//
// Passing an empty string or "-" as the path argument writes to stdout;
// any other value is treated as a file path and written with mode 0600.
type _ = struct{} // blank type to satisfy the file-must-have-non-comment-code rule
