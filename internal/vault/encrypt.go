package vault

import (
	"fmt"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/pkg/vault_logic"
)

func EncryptFile(inputFile, outputFile string, rawKey []byte) error {
	helpers.LogVerbose("Checking if input file already is encrypted")
	encrypted, err := CheckIsEncrypted(inputFile)
	if err == nil && encrypted {
		return fmt.Errorf("input file is already encrypted. Please choose a different input file")
	}

	helpers.LogVerbose("Checking if file exists")
	data, err := ReadEncryptedFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	helpers.LogVerbose("Confirming output file")
	if outputFile == "" {
		outputFile = inputFile + ".vault"
	}
	if helpers.CheckSuffix(inputFile, ".vault") {
		outputFile = inputFile
	}

	finalData, err := EncryptData(rawKey, data)
	if err != nil {
		return err
	}

	err = WriteEncryptedFile(outputFile, finalData)
	if err != nil {
		return err
	}

	return nil
}

func EncryptData(key, data []byte) ([]byte, error) {
	helpers.LogVerbose("Encrypting data")
	finalData, err := vault_logic.EncryptData(key, data)
	if err != nil {
		return nil, fmt.Errorf("error encrypting file: %w", err)
	}

	return finalData, nil
}

func WriteEncryptedFile(outputFile string, data []byte) error {
	helpers.LogVerbose("Writing encrypted data to file")
	err := helpers.WriteFile(outputFile, data)
	if err != nil {
		return fmt.Errorf("error writing encrypted data to file: %w", err)
	}

	return nil
}
