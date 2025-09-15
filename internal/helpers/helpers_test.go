package helpers

import (
	"fmt"
	"os"
	"testing"

	testutils "github.com/ZonCen/Cloak/internal/testUtils"
)

func TestOpenFile(t *testing.T) {
	tmpFile, err := testutils.CreateTempFile()
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Println("Failed to remove temp file:", err)
		}
	}()
	defer func() {
		if err := tmpFile.Close(); err != nil {
			fmt.Println("Failed to close temp file:", err)
		}
	}()

	content := []byte("Hello, World!")

	tests := []struct {
		name      string
		path      string
		wantErr   bool
		wantBytes []byte
	}{
		{"file exists", tmpFile.Name(), false, content},
		{"file does not exist", "non_existent_file.txt", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("OpenFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != string(tt.wantBytes) {
				t.Errorf("OpenFile() = %s, want %s", got, tt.wantBytes)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	content := []byte("Hello World!")
	path := "test_output.txt"
	defer func() {
		if err := os.Remove(path); err != nil {
			fmt.Printf("Failed to remove test file: %v\n", err)
		}
	}()

	if err := WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("WriteFile() wrote %s, want %s", data, content)
	}
}

func TestReadFile(t *testing.T) {
	tmpFile, err := testutils.CreateTempFile()
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	content := []byte("Hello, World!")

	tests := []struct {
		name      string
		path      string
		wantErr   bool
		wantBytes []byte
	}{
		{"file exists", tmpFile.Name(), false, content},
		{"file does not exist", "non_existent_file.txt", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != string(tt.wantBytes) {
				t.Errorf("ReadFile() = %s, want %s", got, tt.wantBytes)
			}
		})
	}
}

func TestCreateFolderIfNotExist(t *testing.T) {
	path := "test_folder/subfolder"
	defer func() {
		if err := os.RemoveAll("test_folder"); err != nil {
			fmt.Printf("Failed to remove test folder: %v\n", err)
		}
	}()

	if err := CreateFolderIfNotExist(path, 0755); err != nil {
		t.Fatalf("CreateFolderIfNotExist() error = %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Folder %s was not created", path)
	}
}

func TestCheckSuffix(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		suffix string
		want   bool
	}{
		{"correct suffix", "file.vault", ".vault", true},
		{"incorrect suffix", "file.txt", ".vault", false},
		{"no suffix", "file", ".vault", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSuffix(tt.path, tt.suffix); got != tt.want {
				t.Errorf("CheckSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveSuffix(t *testing.T) {
	s := "file.vault"
	suffix := ".vault"
	want := "file"
	if got := RemoveSuffix(s, suffix); got != want {
		t.Errorf("RemoveSuffix() = %v, want %v", got, want)
	}
}
