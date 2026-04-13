// Package envfile provides utilities for reading, writing, encrypting,
// decrypting, diffing, merging, auditing, and validating .env files.
//
// # File Format
//
// Each non-blank, non-comment line must follow the pattern:
//
//	KEY=VALUE
//
// Lines beginning with '#' are treated as comments and are preserved
// during round-trip writes where possible.
//
// # Validation
//
// Use [Validate] or [ValidateFile] to check a set of entries for common
// problems such as empty keys, keys that do not conform to the POSIX
// identifier convention ([A-Za-z_][A-Za-z0-9_]*), and duplicate keys.
// [ValidateKey] can be used to check a single key string.
//
// # Encryption
//
// [EncryptValues] and [DecryptValues] operate on []Entry slices, leaving
// keys in plain text while encrypting or decrypting values using AES-GCM
// via the internal/crypto package.
//
// # Diff & Merge
//
// [Diff] compares two []Entry slices and returns a [DiffResult] describing
// added, removed, and changed keys. [Merge] combines a base and overlay
// slice with configurable conflict strategies.
//
// # Audit
//
// [LoadAuditLog] and [Record] provide an append-only JSONL audit trail for
// tracking encryption and rotation events.
package envfile
