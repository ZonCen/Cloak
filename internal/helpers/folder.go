package helpers

import (
	"fmt"
	"os"
)

func CreateFolderIfNotExist(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("failed to create folder %s: %w", path, err)
	}

	return nil
}
