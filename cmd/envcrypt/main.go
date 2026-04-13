// Package main is the entry point for the envcrypt CLI tool.
// It wires together the keystore, envfile, and crypto packages
// to provide commands for encrypting, decrypting, and managing
// per-environment keys for .env files.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envcrypt/internal/envfile"
	"github.com/yourusername/envcrypt/internal/keystore"
)

const defaultKeystorePath = ".envcrypt_keys"

func usage() {
	fmt.Fprintf(os.Stderr, `envcrypt — encrypt and manage .env files with per-environment key rotation

Usage:
  envcrypt <command> [flags]

Commands:
  init      <env>          Generate and store a new key for the given environment
  rotate    <env>          Rotate the key for the given environment
  encrypt   <env> <file>   Encrypt all values in a .env file for the given environment
  decrypt   <env> <file>   Decrypt all values in an encrypted .env file

Flags:
`)
	flag.PrintDefaults()
}

func main() {
	keystorePath := flag.String("keystore", defaultKeystorePath, "path to the keystore file")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "init":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: init requires an <env> argument")
			os.Exit(1)
		}
		runInit(*keystorePath, args[1])

	case "rotate":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: rotate requires an <env> argument")
			os.Exit(1)
		}
		runRotate(*keystorePath, args[1])

	case "encrypt":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "error: encrypt requires <env> and <file> arguments")
			os.Exit(1)
		}
		runEncrypt(*keystorePath, args[1], args[2])

	case "decrypt":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "error: decrypt requires <env> and <file> arguments")
			os.Exit(1)
		}
		runDecrypt(*keystorePath, args[1], args[2])

	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n", command)
		usage()
		os.Exit(1)
	}
}

func runInit(storePath, env string) {
	if err := keystore.GenerateAndStore(storePath, env); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Key generated for environment %q and saved to %s\n", env, storePath)
}

func runRotate(storePath, env string) {
	if err := keystore.Rotate(storePath, env); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Key rotated for environment %q\n", env)
}

func runEncrypt(storePath, env, filePath string) {
	ks, err := keystore.Load(storePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading keystore: %v\n", err)
		os.Exit(1)
	}
	key, ok := ks.Get(env)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: no key found for environment %q — run 'envcrypt init %s' first\n", env, env)
		os.Exit(1)
	}

	pairs, err := envfile.Parse(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env file: %v\n", err)
		os.Exit(1)
	}

	encrypted, err := envfile.EncryptValues(pairs, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encrypting values: %v\n", err)
		os.Exit(1)
	}

	out := encryptedFilePath(filePath, env)
	if err := envfile.Write(out, encrypted); err != nil {
		fmt.Fprintf(os.Stderr, "error writing encrypted file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Encrypted %d values → %s\n", len(encrypted), out)
}

func runDecrypt(storePath, env, filePath string) {
	ks, err := keystore.Load(storePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading keystore: %v\n", err)
		os.Exit(1)
	}
	key, ok := ks.Get(env)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: no key found for environment %q\n", env)
		os.Exit(1)
	}

	pairs, err := envfile.Parse(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env file: %v\n", err)
		os.Exit(1)
	}

	decrypted, err := envfile.DecryptValues(pairs, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decrypting values: %v\n", err)
		os.Exit(1)
	}

	out := decryptedFilePath(filePath)
	if err := envfile.Write(out, decrypted); err != nil {
		fmt.Fprintf(os.Stderr, "error writing decrypted file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Decrypted %d values → %s\n", len(decrypted), out)
}

// encryptedFilePath derives the output path for an encrypted file,
// e.g. ".env" + "production" → ".env.production.enc"
func encryptedFilePath(original, env string) string {
	dir := filepath.Dir(original)
	base := filepath.Base(original)
	return filepath.Join(dir, base+"."+env+".enc")
}

// decryptedFilePath derives the output path for a decrypted file,
// appending ".dec" to avoid overwriting the source.
func decryptedFilePath(original string) string {
	return original + ".dec"
}
