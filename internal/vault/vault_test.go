package vault

import (
	"fmt"
	"os"
	"testing"

	testutils "github.com/ZonCen/Cloak/internal/testUtils"
)

func TestEncryptFile(t *testing.T) {
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

	tests := []struct {
		name    string
		path    string
		key     []byte
		wantErr bool
	}{
		{"valid decryption", tmpFile.Name(), []byte("12345678901234567890123456789012"), false},
		{"invalid key length", tmpFile.Name(), []byte("shortkey"), true},
		{"file does not exist", "non_existent_file.txt", []byte("12345678901234567890123456789012"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncryptFile(tt.key, string(tt.path), string(tt.path)+".vault")
			if err != nil && !tt.wantErr {
				t.Fatalf("EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptFile(t *testing.T) {
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

	key := []byte("12345678901234567890123456789012")
	err = EncryptFile(key, tmpFile.Name(), tmpFile.Name()+".vault")
	if err != nil {
		t.Fatalf("Failed to encrypt file for testing: %v", err)
	}

	content := []byte("Hello, World!")

	tests := []struct {
		name      string
		path      string
		key       []byte
		wantErr   bool
		wantBytes []byte
	}{
		{"valid decryption", tmpFile.Name() + ".vault", key, false, content},
		{"invalid key length", tmpFile.Name() + ".vault", []byte("shortkey"), true, nil},
		{"file does not exist", "non_existent_file.vault", key, true, nil},
		{"invalid file format", tmpFile.Name(), key, true, nil},
		{"invalid key", tmpFile.Name() + ".vault", []byte("12345678901234567890123456789013"), true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecryptFile(tt.key, tt.path, tmpFile.Name())
			if err != nil && !tt.wantErr {
				t.Fatalf("DecryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				got, err := os.ReadFile(tmpFile.Name())
				if err != nil {
					t.Fatalf("Failed to read decrypted file: %v", err)
				}
				if string(got) != string(tt.wantBytes) {
					t.Errorf("DecryptFile() = %s, want %s", got, tt.wantBytes)
				}
			}
		})
	}
}

func TestEncodeAndDecodeKey(t *testing.T) {
	key := []byte("12345678901234567890123456789012")

	encrypted := EncodeKey(key)
	decrypted, err := DecodeKey(encrypted)
	if err != nil {
		t.Fatalf("DecodeKey() error = %v", err)
	}
	if string(decrypted) != string(key) {
		t.Errorf("DecodeKey() = %s, want %s", decrypted, key)
	}
}

func TestIsEncryptedFile(t *testing.T) {
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

	key := []byte("12345678901234567890123456789012")
	err = EncryptFile(key, tmpFile.Name(), tmpFile.Name()+".vault")
	if err != nil {
		t.Fatalf("Failed to encrypt file for testing: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{"encrypted file", tmpFile.Name() + ".vault", true, false},
		{"unencrypted file", tmpFile.Name(), false, false},
		{"non-existent file", "non_existent_file.vault", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsEncryptedFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("IsEncryptedFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("IsEncryptedFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
