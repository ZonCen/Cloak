package vault

import (
	"fmt"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/pkg/vault_logic"
)

func DecryptFile(inputFile, outputFile string, rawKey []byte) error {
	_, err := CheckIsEncrypted(inputFile)
	if err != nil {
		return fmt.Errorf("error while checking if file is encrypted: %w", err)
	}

	data, err := ReadEncryptedFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	helpers.LogVerbose("Confirming output file")
	if outputFile == "" {
		outputFile = inputFile
	}

	decrypted_data, err := DecryptData(rawKey, data)
	if err != nil {
		return err
	}

	err = WriteDecryptedFile(outputFile, decrypted_data)
	if err != nil {
		return err
	}

	return nil
}

func ReadEncryptedFile(inputFile string) ([]byte, error) {
	helpers.LogVerbose("Checking if file exists and is a vault file")
	data, err := helpers.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CheckIsEncrypted(inputFile string) (bool, error) {
	helpers.LogVerbose("Checking if input file is encrypted")
	encrypted, err := vault_logic.IsEncryptedFile(inputFile)
	if err != nil {
		return false, err
	}
	if !encrypted {
		return false, fmt.Errorf("input file is not an encrypted vault file")
	}

	return true, nil
}

func DecryptData(key, data []byte) ([]byte, error) {
	helpers.LogVerbose("Decrypting data")
	decrypted_data, err := vault_logic.DecryptData(key, data)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return decrypted_data, nil
}

func WriteDecryptedFile(outputFile string, data []byte) error {
	helpers.LogVerbose("Writing decrypted data to file")
	err := helpers.WriteFile(outputFile, data)
	if err != nil {
		return fmt.Errorf("error writing decrypted data to file: %w", err)
	}

	return nil
}
