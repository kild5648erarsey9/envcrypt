// Package envfile handles reading and writing .env files and provides
// helpers to encrypt or decrypt their values using AES-GCM keys supplied
// by the keystore package.
//
// Typical workflow:
//
//	// Encrypt an existing plaintext .env file:
//	 ef, err := envfile.Parse(".env")
//	 if err != nil { ... }
//	 if err := envfile.EncryptValues(ef, key); err != nil { ... }
//	 if err := envfile.Write(".env.enc", ef); err != nil { ... }
//
//	// Decrypt back to plaintext:
//	 ef, err := envfile.Parse(".env.enc")
//	 if err != nil { ... }
//	 if err := envfile.DecryptValues(ef, key); err != nil { ... }
//	 if err := envfile.Write(".env", ef); err != nil { ... }
package envfile
