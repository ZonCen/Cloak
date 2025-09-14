package vault

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomByteKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

func EncodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
