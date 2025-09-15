package vault

import (
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/ZonCen/Cloak/internal/helpers"
)

func IsEncryptedFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	header := make([]byte, 5)
	n, err := f.Read(header)
	if err != nil {
		return false, fmt.Errorf("failed to read file header: %w", err)
	}

	if n < 5 {
		return false, fmt.Errorf("file too short to be a valid Cloak encrypted file")
	}

	if string(header[:4]) == "CLOK" && header[4] == 1 {
		return true, nil
	}

	return false, nil
}

func ExtractNonce(data []byte, gcm cipher.AEAD) ([]byte, []byte, error) {
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, nil, fmt.Errorf("ciphertext too short: %d bytes", len(data))
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]

	return nonce, cipherText, nil
}

func DecryptFile(key []byte, inputPath, outputPath string) error {
	plaintext, err := ReadEncryptedFile(key, inputPath)
	if err != nil {
		return err
	}

	return helpers.WriteFile(outputPath, plaintext)
}

func ReadEncryptedFile(key []byte, inputPath string) ([]byte, error) {
	data, err := helpers.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	block, err := CreateCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := GenerateGCM(block)
	if err != nil {
		return nil, err
	}

	header, nonce, cipherText, err := ExtractHeaderAndNonce(data, gcm)
	if err != nil {
		return nil, err
	}

	if err := ValidateHeader(header); err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

func GetEnv(env string) ([]byte, error) {
	helpers.LogVerbose("Fecthing key from environment variable")
	key := os.Getenv(env)
	if key == "" {
		return nil, fmt.Errorf("%s environment variable is not set", env)
	}
	rawKey, err := DecodeKey(key)
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}
	if len(rawKey) != 32 {
		return nil, fmt.Errorf("invalid key length: must be 32 bytes for AES-256")
	}

	return rawKey, nil
}

func ExtractHeaderAndNonce(data []byte, gcm cipher.AEAD) ([]byte, []byte, []byte, error) {
	helpers.LogVerbose("Extracting header and nonce from data")
	headerSize := 5
	nonceSize := gcm.NonceSize()
	if len(data) < headerSize+nonceSize {
		return nil, nil, nil, fmt.Errorf("file too short: got %d bytes", len(data))
	}

	header := data[:headerSize]
	nonce := data[headerSize : headerSize+nonceSize]
	ciphertext := data[headerSize+nonceSize:]

	return header, nonce, ciphertext, nil
}

func ValidateHeader(header []byte) error {
	helpers.LogVerbose("Validating file header")
	if string(header[:4]) != "CLOK" {
		return fmt.Errorf("invalid file format: missing CLOK magic bytes")
	}
	version := header[4]
	if version != 1 {
		return fmt.Errorf("unsupported file version: %d", version)
	}
	return nil
}
