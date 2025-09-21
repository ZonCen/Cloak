package vault_logic

import (
	"fmt"
	"os"
)

type FileHeader struct {
	Magic   [4]byte
	Version byte
}

func EncryptData(key, data []byte) ([]byte, error) {
	block, err := CreateBlock(key)
	if err != nil {
		return nil, err
	}

	gcm, err := CreateCipher(block)
	if err != nil {
		return nil, err
	}

	nonce, err := GenerateNonce(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	header := BuildHeader()
	finalData := append(header, nonce...)
	finalData = append(finalData, ciphertext...)

	return finalData, nil
}

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

func ReadEncryptedData(key, data []byte) ([]byte, error) {
	block, err := CreateBlock(key)
	if err != nil {
		return nil, err
	}

	gcm, err := CreateCipher(block)
	if err != nil {
		return nil, err
	}

	header, nonce, cipherText, err := GenerateHeaderAndNonce(data, gcm)
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
