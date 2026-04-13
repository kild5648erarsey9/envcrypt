# envcrypt

> A CLI tool to encrypt and manage `.env` files with per-environment key rotation support.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envcrypt.git
cd envcrypt && go build -o envcrypt .
```

---

## Usage

**Encrypt a `.env` file:**

```bash
envcrypt encrypt --env production --file .env.production
```

**Decrypt a `.env` file:**

```bash
envcrypt decrypt --env production --file .env.production.enc
```

**Rotate the key for an environment:**

```bash
envcrypt rotate --env production
```

**List managed environments:**

```bash
envcrypt list
```

Keys are stored locally in `~/.envcrypt/keys/` and are scoped per environment, allowing independent rotation without affecting other environments.

---

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--env` | Target environment name | `development` |
| `--file` | Path to the `.env` file | `.env` |
| `--output` | Output file path | `<file>.enc` |

---

## License

This project is licensed under the [MIT License](LICENSE).