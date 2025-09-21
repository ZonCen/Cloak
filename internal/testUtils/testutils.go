package testutils

import (
	"os"
	"testing"
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

func CreateTempFileWithCleanup(t *testing.T) *os.File {
	t.Helper()
	f, err := CreateTempFile()
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		err = f.Close()
		if err != nil {
			t.Logf("Failed to close file %v, with error: %v", f, err)
		}
		err = os.Remove(f.Name())
		if err != nil {
			t.Logf("Failed to remove file %v, with error: %v", f, err)
		}
	})

	return f
}
