package vault

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/pkg/vault_logic"
)

func GetEnv(env string) (string, error) {
	helpers.LogVerbose("Fecthing key from environment variable")
	key := os.Getenv(env)
	if key == "" {
		return "", fmt.Errorf("%s environment variable is not set", env)
	}
	return key, nil
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	return editor
}

func CreateTempFile() (*os.File, error) {
	return os.CreateTemp("", "cloak_edit_*.tmp")
}

func LaunchEditor(editor, filePath string) error {
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch editor: %w", err)
	}

	return nil
}

func GetKey() ([]byte, error) {
	helpers.LogVerbose("Getting envoronment variable")
	key, err := GetEnv("CLOAK_KEY")
	if err != nil {
		return nil, err
	}

	rawKey, err := vault_logic.DecodeKey(key)
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}
	if len(rawKey) != 32 {
		return nil, fmt.Errorf("invalid key length: must be 32 bytes for AES-256")
	}

	return rawKey, nil
}
