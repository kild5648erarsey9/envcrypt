// Package keystore manages per-environment AES encryption keys for envcrypt.
//
// Keys are persisted as a JSON file on disk (mode 0600) and indexed by
// environment name (e.g. "production", "staging", "development").
//
// Typical usage:
//
//	ks, err := keystore.Load(".envcrypt/keys.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Generate a key for a new environment.
//	hexKey, err := keystore.GenerateAndStore(ks, "production")
//
//	// Rotate an existing key.
//	result, err := keystore.Rotate(ks, "production")
//	// result.OldKey and result.NewKey can be used to re-encrypt .env data.
//
//	if err := ks.Save(); err != nil {
//		log.Fatal(err)
//	}
//
// The keystore file should be kept outside of version control. envcrypt's
// default .gitignore excludes .envcrypt/keys.json automatically.
package keystore
