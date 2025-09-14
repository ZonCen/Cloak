package vault

import (
	"fmt"
	"os"
	"os/exec"
)

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
