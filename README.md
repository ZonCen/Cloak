# Cloak

![Cloak Logo](https://via.placeholder.com/150)

Cloak is an **open-source encryption CLI tool** for securely managing and storing sensitive files. It provides a lightweight, extensible way to encrypt and decrypt files, making it safe to store configuration files, secrets, or other sensitive data in version control. Currently, Cloak focuses on the CLI functionality, including file encryption, decryption, editing, and key initialization.

---

## Features

- **AES-256-GCM encryption** for authenticated, secure file storage  
- **File metadata header** to store magic bytes and version info  
- **CLI command to initialize encryption key** (`init`)  
- **Command-line editor integration** to safely edit encrypted files in your default editor  
- **Cross-platform CLI** built with [Cobra](https://github.com/spf13/cobra)  

---

## Installation

Clone the repository:

```bash
git clone https://github.com/ZonCen/Cloak.git
cd Cloak
```

Build the CLI:

```bash
go build -o cloak ./cmd
```

Make sure your `$GOPATH/bin` is in your PATH if using `go install`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Getting Started

### 1. Initialize Cloak

```bash
cloak init
```

- Generates a new 32-byte encryption key (`CLOAK_KEY`) and stores it securely.  
- You **only need to export `CLOAK_KEY`** into your environment after running `init` for the first session:

```bash
export CLOAK_KEY="<generated-key-from-init>"
```

### 2. Encrypt a file

```bash
cloak encrypt config.yaml config.yaml.vault
```

- Encrypts `config.yaml` and saves as `config.yaml.vault`.

### 3. Decrypt a file

```bash
cloak decrypt config.yaml.vault config.yaml
```

- Decrypts `config.yaml.vault` and writes plaintext to `config.yaml`.

### 4. Edit an encrypted file

```bash
cloak edit config.yaml.vault
```

- Decrypts the file into a temporary buffer  
- Opens your default editor (`$EDITOR` or `vi`)  
- Re-encrypts the file after editing  

---

## File Structure

- `cmd` – CLI commands (`init`, `encrypt`, `decrypt`, `edit`)  
- `internal/helpers` – file and logging helpers  
- `tests` – integration and end-to-end tests  

---

## Security Notes

- Uses **AES-256-GCM** with unique nonce for each file.  
- File header contains `CLOK` magic bytes and version info for validation.  
- Keys should be **stored securely** in environment variables or other secret management systems.  
- Always backup your original `.vault` files before editing.

---

## Contributing

Contributions are welcome!  

1. Fork the repository  
2. Create a feature branch (`git checkout -b feature-name`)  
3. Write tests for new functionality  
4. Submit a pull request  

---

## License

MIT License – see [LICENSE](LICENSE) for details.  

---

## Diagram: Cloak CLI Workflow

```text
+----------------+      encrypt       +----------------+
| Plaintext File | ----------------> | Encrypted File |
|   config.yaml  |                    | config.yaml.vault |
+----------------+                     +----------------+
          ^                                      |
          |                                      |
          |         decrypt / edit               |
          |                                      v
+----------------+ <----------------- +----------------+
| Plaintext File |                    | Temporary File |
|   config.yaml  |                    |  opened in editor |
+----------------+                     +----------------+
```

This diagram illustrates the **round-trip workflow** for the CLI: plaintext file → encrypted `.vault` file → temp file → editor → re-encrypt → original `.vault` file.

