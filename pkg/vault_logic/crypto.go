package vault_logic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func CreateBlock(key []byte) (cipher.Block, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return block, nil
}

func CreateCipher(block cipher.Block) (cipher.AEAD, error) {
	if block == nil {
		return nil, fmt.Errorf("block cipher is nil")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	return gcm, nil
}

func GenerateNonce(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size %d", size)
	}
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return nonce, nil
}

func GenerateRandomByteKey(length int) ([]byte, error) {
	key, err := GenerateNonce(length)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateHeaderAndNonce(data []byte, gcm cipher.AEAD) ([]byte, []byte, []byte, error) {
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
