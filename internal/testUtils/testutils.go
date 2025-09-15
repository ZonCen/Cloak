package testutils

import (
	"os"
)

func CreateTempFile() (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		return nil, err
	}

	content := []byte("Hello, World!")
	if _, err := tmpFile.Write(content); err != nil {
		return nil, err
	}

	return tmpFile, nil
}
