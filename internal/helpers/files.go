package helpers

import (
	"fmt"
	"io"
	"os"
)

func OpenFile(path string) ([]byte, error) {
	infile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() {
		if cerr := infile.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", path, cerr)
		}
	}()

	info, err := io.ReadAll(infile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return info, nil
}

func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return data, nil
}
