package vault_logic

func DecryptData(key, data []byte) ([]byte, error) {
	plaintext, err := ReadEncryptedData(key, data)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
