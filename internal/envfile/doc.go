// Package envfile provides utilities for working with .env files in the
// envcrypt tool.
//
// # Parsing and Writing
//
// Parse reads a .env file from disk into an ordered slice of key-value pairs,
// preserving blank lines and comments. Write serialises that slice back to
// disk in a deterministic, human-readable format.
//
// # Encryption and Decryption
//
// EncryptValues encrypts every value in an env-file map using the AES-GCM
// primitives from internal/crypto, base64-encoding each ciphertext so the
// result remains a valid .env file. DecryptValues reverses the process.
//
// # Diffing
//
// Diff compares two env-file maps and returns a structured summary of which
// keys were added, removed, or changed between them — useful for auditing
// changes before committing an encrypted file.
//
// # Audit Logging
//
// Record, LoadAuditLog, and AuditLog provide a lightweight append-only JSON
// audit trail. Every encrypt, decrypt, and rotate operation performed by the
// CLI appends an AuditEvent (timestamp, environment, operation, affected keys)
// to a local audit.json file so operators can review what changed and when.
package envfile
