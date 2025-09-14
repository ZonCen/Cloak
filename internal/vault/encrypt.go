package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/ZonCen/Cloak/internal/helpers"
)

type FileHeader struct {
	Magic   [4]byte
	Version byte
}

func EncryptFile(key []byte, inputPath, outputPath string) error {
	plaintext, err := helpers.OpenFile(inputPath)
	if err != nil {
		return err
	}
	block, err := CreateCipher(key)
	if err != nil {
		return err
	}

	gcm, err := GenerateGCM(block)
	if err != nil {
		return err
	}

	nonce, err := GenerateNonce(gcm.NonceSize())
	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	header := BuildHeader()
	finalData := append(header, nonce...)
	finalData = append(finalData, ciphertext...)

	return helpers.WriteFile(outputPath, finalData)
}

func CreateCipher(key []byte) (cipher.Block, error) {
	helpers.LogVerbose("Creating AES cipher")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return block, nil
}

func GenerateGCM(block cipher.Block) (cipher.AEAD, error) {
	helpers.LogVerbose("Generating GCM")
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	return gcm, nil
}

func GenerateNonce(size int) ([]byte, error) {
	helpers.LogVerbose("Generating nonce")
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return nonce, nil
}

func BuildHeader() []byte {
	helpers.LogVerbose("Building file header")
	var header FileHeader
	copy(header.Magic[:], []byte("CLOK"))
	header.Version = 1

	buf := make([]byte, 0, 5)
	buf = append(buf, header.Magic[:]...)
	buf = append(buf, header.Version)

	return buf
}
