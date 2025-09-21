package vault

import (
	"os"
	"testing"

	testutils "github.com/ZonCen/Cloak/internal/testUtils"
)

var (
	rawKey   = []byte("12345678901234567890123456789012")
	testData = []byte("Hello, World!")
)

func TestEncryptFile(t *testing.T) {
	tmpFile := testutils.CreateTempFileWithCleanup(t)

	tests := []struct {
		name    string
		path    string
		key     []byte
		wantErr bool
	}{
		{"valid decryption", tmpFile.Name(), rawKey, false},
		{"invalid key length", tmpFile.Name(), []byte("shortkey"), true},
		{"file does not exist", "non_existent_file.txt", rawKey, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncryptFile(string(tt.path), string(tt.path)+".vault", tt.key)
			if err != nil && !tt.wantErr {
				t.Fatalf("EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		key     []byte
		wantErr bool
	}{
		{"Working encrypt data", testData, rawKey, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedData, err := EncryptData(tt.key, tt.data)
			if err != nil {
				t.Fatalf("EncryptData() error = %v", err)
			}
			if len(encryptedData) <= len(tt.data) {
				t.Errorf("EncryptData() = %v, want length greater than %d", encryptedData, len(tt.data))
			}
		})
	}
}

func TestWriteEncryptedFile(t *testing.T) {
	tmpFile := testutils.CreateTempFileWithCleanup(t)

	tests := []struct {
		name    string
		data    []byte
		key     []byte
		wantErr bool
	}{
		{"Working to write Encrypted file", testData, rawKey, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := EncryptData(tt.key, tt.data)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to encrypt data: %v", err)
			}

			err = WriteEncryptedFile(tmpFile.Name(), data)
			if err != nil && !tt.wantErr {
				t.Fatalf("WriteEncryptedFile() error = %v", err)
			}
		})
	}
}

func TestDecryptFile(t *testing.T) {
	tmpFile := testutils.CreateTempFileWithCleanup(t)

	tests := []struct {
		name    string
		path    string
		key     []byte
		wantErr bool
	}{
		{"valid decryption", tmpFile.Name(), rawKey, false},
		{"invalid key length", tmpFile.Name(), []byte("shortkey"), true},
		{"file does not exist", "non_existent_file.txt", rawKey, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncryptFile(tmpFile.Name(), tmpFile.Name()+".vault", tt.key)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to encrypt file: %v", err)
			}
			err = DecryptFile(string(tt.path)+".vault", string(tt.path)+".vault", tt.key)
			if err != nil && !tt.wantErr {
				t.Fatalf("Decryptfile() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadEncryptedfile(t *testing.T) {
	tmpFile := testutils.CreateTempFileWithCleanup(t)

	tests := []struct {
		name    string
		path    string
		wantErr bool
		data    []byte
	}{
		{"can read file", tmpFile.Name(), false, testData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncryptFile(tmpFile.Name(), tmpFile.Name()+".vault", rawKey)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to encrypt file: %v", err)
			}
			data, err := ReadEncryptedFile(tt.path + ".vault")
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to read encrypted file: %v", err)
			}

			if string(data) == string(tt.data) && !tt.wantErr {
				t.Fatalf("Data is not encrypted")
			}
		})
	}
}

func TestCheckIsEncrypted(t *testing.T) {
	tmpfile := testutils.CreateTempFileWithCleanup(t)
	vaultFile := tmpfile.Name() + ".vault"

	t.Cleanup(func() {
		if err := os.Remove(vaultFile); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove vault fule: %v", err)
		}
	})

	err := EncryptFile(tmpfile.Name(), vaultFile, rawKey)
	if err != nil {
		t.Fatalf("Failed to create encrypted temp file %v", vaultFile)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Is encrypted", vaultFile, false},
		{"Is not encrypted", tmpfile.Name(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := CheckIsEncrypted(tt.path)
			if err != nil && !tt.wantErr {
				t.Fatalf("Error when checking if file is encrypted: %v", err)
			}
			if encrypted == tt.wantErr {
				t.Fatalf("CheckIsEncrypted() failed, expected %v, got %v", tt.wantErr, encrypted)
			}
		})
	}
}

func TestDecryptData(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		data    []byte
		wantErr bool
	}{
		{
			name:    "example",
			key:     rawKey,
			data:    testData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptData(tt.key, tt.data)
			if err != nil {
				t.Fatalf("Error when encrypting data: %v", err)
			}
			if string(encrypted) == string(tt.data) {
				t.Fatal("The encryption have not encrypted the data")
			}

			decrypted_data, err := DecryptData(tt.key, encrypted)
			if err != nil && !tt.wantErr {
				t.Fatalf("Error when decrypting data: %v", err)
			}

			if string(decrypted_data) != string(tt.data) {
				t.Fatalf("DecryptData() = %v, wanted same text as %v", decrypted_data, tt.data)
			}
		})
	}
}

func TestReadEcnryptedFile(t *testing.T) {
	tmpFile := testutils.CreateTempFileWithCleanup(t)
	vaultFile := tmpFile.Name() + ".vault"
	tmpFile2 := testutils.CreateTempFileWithCleanup(t)

	err := EncryptFile(tmpFile.Name(), tmpFile.Name()+".vault", rawKey)
	if err != nil {
		t.Fatalf("Error encrypting file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(vaultFile); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove vault fule: %v", err)
		}
	})

	tests := []struct {
		name    string
		path    string
		wantErr bool
		data    []byte
	}{
		{
			name:    "File exist and is encrypted",
			path:    vaultFile,
			wantErr: false,
			data:    testData,
		},
		{
			name:    "File exist but is not encrypted",
			path:    tmpFile2.Name() + ".vault",
			wantErr: true,
			data:    testData,
		},
		{
			name:    "File does not exist",
			path:    "some_non_existing_file.vault",
			wantErr: true,
			data:    testData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ReadEncryptedFile(tt.path)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to read the encrypted file %v, error: %v", tt.path, err)
			}

			if len(data) <= len(tt.data) && !tt.wantErr {
				t.Fatalf("ReadEncryptedFile() = %v, want length greater than %d", data, len(tt.data))
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name    string
		env_key string
		env_val string
		wantErr bool
	}{
		{
			name:    "Getting existing environment variable",
			env_key: "CLOAK_KEY_TEST",
			env_val: string(rawKey),
			wantErr: false,
		},
		{
			name:    "Getting none existing envornment variable",
			env_key: "CLOAKY_DOAKY_MOAKY",
			env_val: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := os.Setenv(tt.env_key, ""); err != nil {
					t.Logf("Failed to reset environment %v, back to %v", tt.env_key, "")
				}
			}()
			err := os.Setenv(tt.env_key, tt.env_val)
			if err != nil {
				t.Fatalf("Failed to set environment variable: %v with error: %v", tt.env_key, err)
			}
			envVar, err := GetEnv(tt.env_key)
			if err != nil && !tt.wantErr {
				t.Fatalf("Error getting the environment variable: %v", err)
			}

			if envVar != tt.env_val {
				t.Fatalf("GetEnv() = %v, expected %v", envVar, tt.env_val)
			}
		})
	}
}

func TestGetEditor(t *testing.T) {
	tests := []struct {
		name    string
		env_key string
		env_val string
		expRes  string
		wantErr bool
	}{
		{
			name:    "Getting EDITOR=vim",
			env_key: "EDITOR",
			env_val: "vim",
			expRes:  "vim",
		},
		{
			name:    "Getting EDITOR=''",
			env_key: "EDITOR",
			env_val: "",
			expRes:  "vi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVal := os.Getenv(tt.env_key)
			defer func() {
				if err := os.Setenv(tt.env_key, oldVal); err != nil {
					t.Logf("Failed to reset environment %v, back to %v", tt.env_key, oldVal)
				}
			}()
			err := os.Setenv(tt.env_key, tt.env_val)
			if err != nil {
				t.Fatalf("Error setting the environment variable: %v", err)
			}

			editor := GetEditor()

			if editor != tt.expRes {
				t.Fatalf("GetEditor() = %v, expected %v", editor, tt.expRes)
			}
		})
	}
}
